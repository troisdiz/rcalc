package rcalc

import (
	"fmt"
	"strings"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

func MarshallActions(registry *ActionRegistry, actions []Action) ([]*protostack.Action, error) {
	var protoActions []*protostack.Action
	for _, rAction := range actions {
		protoAction, err := rAction.MarshallFunc()(registry, rAction)
		if err != nil {
			return nil, err
		}
		protoActions = append(protoActions, protoAction)
	}
	return protoActions, nil
}

func UnMarshallActions(registry *ActionRegistry, protoActions []*protostack.Action) ([]Action, error) {
	var actions []Action
	for _, action := range protoActions {
		loopAction, err := registry.CreateActionFromProto(action)
		if err != nil {
			return nil, err
		}
		actions = append(actions, loopAction)
	}

	return actions, nil
}

/*
	Action used to add a variable to the stack in a program (can be seen as Variable wrapper)
*/
type VariablePutOnStackActionDesc struct {
	value Variable
}

// VariablePutOnStackActionDesc implements Action
var _ Action = (*VariablePutOnStackActionDesc)(nil)

func (a *VariablePutOnStackActionDesc) NbArgs() int {
	return 0
}

func (a *VariablePutOnStackActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (a *VariablePutOnStackActionDesc) Apply(runtimeContext *RuntimeContext) error {
	runtimeContext.stack.Push(a.value)
	return nil
}

func (a *VariablePutOnStackActionDesc) OpCode() string {
	return "__hidden__" + "PutOnStack"
}

func (a *VariablePutOnStackActionDesc) String() string {
	return fmt.Sprintf("%s(%s)", a.OpCode(), a.value.String())
}

func (a *VariablePutOnStackActionDesc) Display() string {
	return a.value.display()
}

func (a *VariablePutOnStackActionDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {
		variablePutOnStackActionDesc := action.(*VariablePutOnStackActionDesc)

		protoVariable, err := CreateProtoFromVariable(variablePutOnStackActionDesc.value)
		if err != nil {
			return nil, err
		}
		protoPutVariableOnStackAction := &protostack.PutVariableOnStackAction{Value: protoVariable}
		protoAction := &protostack.Action_PutVariableOnStackAction{PutVariableOnStackAction: protoPutVariableOnStackAction}
		return &protostack.Action{
				Type:       protostack.ActionType_PUT_VARIABLE_ON_STACK,
				OpCode:     action.OpCode(),
				RealAction: protoAction},
			nil
	}
}

func (a *VariablePutOnStackActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		protoPutVariableOnStackAction := protoAction.GetPutVariableOnStackAction()
		variableValue, err := CreateVariableFromProto(reg, protoPutVariableOnStackAction.Value)
		if err != nil {
			return nil, err
		}
		return &VariablePutOnStackActionDesc{value: variableValue}, nil
	}
}

/*
	Looks for a local variable named VariableEvaluationDesc.varName in the RuntimeContext and put its value
	on the stack.
*/
type VariableEvaluationActionDesc struct {
	varName string
}

// VariableEvaluationActionDesc implements Action
var _ Action = (*VariableEvaluationActionDesc)(nil)

func (a *VariableEvaluationActionDesc) NbArgs() int {
	return 0
}

func (a *VariableEvaluationActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (a *VariableEvaluationActionDesc) Apply(runtimeContext *RuntimeContext) error {
	value, err := runtimeContext.GetVariableValue(a.varName)
	if err != nil {
		return err
	}
	runtimeContext.stack.Push(value)
	return nil
}

func (a *VariableEvaluationActionDesc) OpCode() string {
	return "__hidden__" + "VariableEvaluation"
}

func (a *VariableEvaluationActionDesc) String() string {
	return fmt.Sprintf("%s(%s)", a.OpCode(), a.varName)
}

func (a *VariableEvaluationActionDesc) Display() string {
	return a.varName
}

func (a *VariableEvaluationActionDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {
		protoEvaluationAction := &protostack.VariableEvaluationAction{VarName: action.(*VariableEvaluationActionDesc).varName}
		return &protostack.Action{Type: protostack.ActionType_VARIABLE_EVALUATION,
				OpCode: action.OpCode(),
				RealAction: &protostack.Action_VariableEvaluationAction{VariableEvaluationAction: protoEvaluationAction,
				},
			},
			nil
	}
}

func (a *VariableEvaluationActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		protoVariableEvaluationAction := protoAction.GetVariableEvaluationAction()
		return &VariableEvaluationActionDesc{varName: protoVariableEvaluationAction.VarName}, nil
	}
}

/*
	Action to execute if then [else] code structures
*/
type IfThenElseActionDesc struct {
	ifActions   []Action
	thenActions []Action
	elseActions []Action
}

var _ Action = (*IfThenElseActionDesc)(nil)

func (a *IfThenElseActionDesc) OpCode() string {
	return "__hidden__" + "IfThenElse"
}

func (a *IfThenElseActionDesc) NbArgs() int {
	return 0
}

func (a *IfThenElseActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (a *IfThenElseActionDesc) Apply(runtimeContext *RuntimeContext) error {
	for _, action := range a.ifActions {
		err := runtimeContext.RunAction(action)
		if err != nil {
			return err
		}
	}
	boolVar, err := runtimeContext.stack.Pop()
	if err != nil {
		return err
	}
	var nextActions []Action
	if boolVar.asBooleanVar().value {
		nextActions = a.thenActions
	} else {
		nextActions = a.elseActions
	}
	for _, action := range nextActions {
		err := runtimeContext.RunAction(action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *IfThenElseActionDesc) Display() string {
	ifActionsStr := []string{}
	for _, action := range a.ifActions {
		ifActionsStr = append(ifActionsStr, action.Display())
	}
	thenActionsStr := []string{}
	for _, action := range a.thenActions {
		thenActionsStr = append(thenActionsStr, action.Display())
	}
	elseActionsStr := []string{}
	for _, action := range a.elseActions {
		elseActionsStr = append(elseActionsStr, action.Display())
	}
	if len(elseActionsStr) == 0 {

		return fmt.Sprintf("if %s then %s end",
			strings.Join(ifActionsStr, " "),
			strings.Join(thenActionsStr, " "))
	} else {
		return fmt.Sprintf("if %s then %s else %s end",
			strings.Join(ifActionsStr, " "),
			strings.Join(thenActionsStr, " "),
			strings.Join(elseActionsStr, " "))
	}
}

func (a *IfThenElseActionDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {
		ifProtoActions, err := MarshallActions(reg, action.(*IfThenElseActionDesc).ifActions)
		if err != nil {
			return nil, err
		}
		thenProtoActions, err := MarshallActions(reg, action.(*IfThenElseActionDesc).thenActions)
		if err != nil {
			return nil, err
		}
		elseProtoActions, err := MarshallActions(reg, action.(*IfThenElseActionDesc).elseActions)
		if err != nil {
			return nil, err
		}
		protoIfThenElseAction := &protostack.IfThenElseAction{
			IfActions:   ifProtoActions,
			ThenActions: thenProtoActions,
			ElseActions: elseProtoActions,
		}
		return &protostack.Action{
			Type:       protostack.ActionType_IF_THEN_ELSE,
			OpCode:     a.OpCode(),
			RealAction: &protostack.Action_IfThenElseAction{IfThenElseAction: protoIfThenElseAction},
		}, nil
	}
}

func (a *IfThenElseActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		ifThenElseAction := protoAction.GetIfThenElseAction()
		ifActions, err := UnMarshallActions(reg, ifThenElseAction.IfActions)
		if err != nil {
			return nil, err
		}
		thenActions, err := UnMarshallActions(reg, ifThenElseAction.ThenActions)
		if err != nil {
			return nil, err
		}
		elseActions, err := UnMarshallActions(reg, ifThenElseAction.ElseActions)
		if err != nil {
			return nil, err
		}
		return &IfThenElseActionDesc{
			ifActions:   ifActions,
			thenActions: thenActions,
			elseActions: elseActions,
		}, nil
	}
}

/*
	Action to execute start ... next loops
*/
type StartNextLoopActionDesc struct {
	actions []Action
}

func (a *StartNextLoopActionDesc) NbArgs() int {
	return 2
}

// StartNextLoopActionDesc  implements Action
var _ Action = (*StartNextLoopActionDesc)(nil)

func (a *StartNextLoopActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	for i := 0; i <= 1; i++ {
		if elts[i].getType() != TYPE_NUMERIC && elts[i].asNumericVar().value.IsInteger() {
			return false, fmt.Errorf("%s at stack level %d is not an integer", elts[i].String(), i+1)
		}
	}
	return true, nil
}

func (a *StartNextLoopActionDesc) Apply(runtimeContext *RuntimeContext) error {
	boundaries, err := runtimeContext.stack.PopN(2)
	if err != nil {
		return err
	}
	start := boundaries[0].asNumericVar().value.IntPart()
	end := boundaries[1].asNumericVar().value.IntPart()
	for i := start; i <= end; i++ {
		for _, action := range a.actions {
			err = action.Apply(runtimeContext)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *StartNextLoopActionDesc) OpCode() string {
	return "__hidden__" + "StartNextLoop"
}

func (a *StartNextLoopActionDesc) String() string {
	return fmt.Sprintf("%s ()", a.OpCode())
}

func (a *StartNextLoopActionDesc) Display() string {
	actionsStr := []string{}
	for _, action := range a.actions {
		actionsStr = append(actionsStr, action.Display())
	}
	return fmt.Sprintf("start %s next", strings.Join(actionsStr, " "))
}

func (a *StartNextLoopActionDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {
		var protoActions []*protostack.Action
		for _, action := range action.(*StartNextLoopActionDesc).actions {
			protoAction, err := action.MarshallFunc()(reg, action)
			if err != nil {
				return nil, err
			}
			protoActions = append(protoActions, protoAction)
		}
		protoStartNextLoopAction := &protostack.StartNextLoopAction{Actions: protoActions}
		return &protostack.Action{
				Type:   protostack.ActionType_START_NEXT,
				OpCode: action.OpCode(),
				RealAction: &protostack.Action_StartNextLoopAction{
					StartNextLoopAction: protoStartNextLoopAction}},
			nil
	}
}

func (a *StartNextLoopActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		protoActions := protoAction.GetStartNextLoopAction().GetActions()
		var actions []Action
		for _, loopProtoAction := range protoActions {
			action, err := reg.CreateActionFromProto(loopProtoAction)
			if err != nil {
				return nil, err
			}
			actions = append(actions, action)
		}
		return &StartNextLoopActionDesc{actions: actions}, nil
	}
}

func CreateStartNextLoopAction(actions []Action) *StartNextLoopActionDesc {
	return &StartNextLoopActionDesc{actions: actions}
}

type ForNextLoopActionDesc struct {
	varName string
	actions []Action
}

// ForNextLoopActionDesc implements Action
var _ Action = (*ForNextLoopActionDesc)(nil)

func (a *ForNextLoopActionDesc) OpCode() string {
	return "__hidden__" + "ForNextLoop"
}

func (a *ForNextLoopActionDesc) NbArgs() int {
	return 2
}

func (a *ForNextLoopActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	for i := 0; i <= 1; i++ {
		if elts[i].getType() != TYPE_NUMERIC && elts[i].asNumericVar().value.IsInteger() {
			return false, fmt.Errorf("%s at stack level %d is not an integer", elts[i].String(), i+1)
		}
	}
	return true, nil
}

func (a *ForNextLoopActionDesc) Apply(runtimeContext *RuntimeContext) error {

	runtimeContext.EnterNewScope()
	defer func() {
		runtimeContext.LeaveScope()
	}()

	boundaries, err := runtimeContext.stack.PopN(2)
	if err != nil {
		return err
	}
	start := boundaries[0].asNumericVar().value.IntPart()
	end := boundaries[1].asNumericVar().value.IntPart()
	// new scope

	for i := start; i <= end; i++ {
		// set var value
		err := runtimeContext.SetVariableValue(a.varName, CreateNumericVariableFromInt(int(i)))
		if err != nil {
			return err
		}
		for _, action := range a.actions {
			err = action.Apply(runtimeContext)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *ForNextLoopActionDesc) Display() string {
	actionsStr := []string{}
	for _, action := range a.actions {
		actionsStr = append(actionsStr, action.Display())
	}
	return fmt.Sprintf("for %s %s next", a.varName, strings.Join(actionsStr, " "))
}

func (a *ForNextLoopActionDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {
		forNextLoopActionDesc := action.(*ForNextLoopActionDesc)
		var protoActions []*protostack.Action
		for _, loopAction := range forNextLoopActionDesc.actions {
			protoLoopAction, err := loopAction.MarshallFunc()(reg, loopAction)
			if err != nil {
				return nil, err
			}
			protoActions = append(protoActions, protoLoopAction)
		}
		protoForNextLoopAction := &protostack.ForNextLoopAction{
			VarName: forNextLoopActionDesc.varName,
			Actions: protoActions,
		}
		return &protostack.Action{
			Type:       protostack.ActionType_FOR_NEXT,
			OpCode:     action.OpCode(),
			RealAction: &protostack.Action_ForNextLoopAction{ForNextLoopAction: protoForNextLoopAction},
		}, nil
	}
}

func (a *ForNextLoopActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		var loopActions []Action
		protoForNextLoopAction := protoAction.GetForNextLoopAction()
		for _, protoLoopAction := range protoForNextLoopAction.GetActions() {
			loopAction, err := reg.CreateActionFromProto(protoLoopAction)
			if err != nil {
				return nil, err
			}
			loopActions = append(loopActions, loopAction)
		}
		return &ForNextLoopActionDesc{
			varName: protoForNextLoopAction.VarName,
			actions: loopActions,
		}, nil
	}
}

type EvalProgramActionDesc struct {
	program *ProgramVariable
}

// EvalProgramActionDesc implements Action
var _ Action = (*EvalProgramActionDesc)(nil)

func (e *EvalProgramActionDesc) Display() string {
	return e.program.display()
}

func (e *EvalProgramActionDesc) OpCode() string {
	return "__hidden__" + "EvalProgram"
}

func (e *EvalProgramActionDesc) NbArgs() int {
	return 0
}

func (e *EvalProgramActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (e *EvalProgramActionDesc) Apply(runtimeContext *RuntimeContext) error {
	return executeProgram(runtimeContext, e.program)
}

func executeProgram(runtimeContext *RuntimeContext, program *ProgramVariable) error {
	runtimeContext.EnterNewScope()
	defer func() { runtimeContext.LeaveScope() }()

	for _, action := range program.actions {
		err := runtimeContext.RunAction(action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EvalProgramActionDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {
		evalProgActionDesc := action.(*EvalProgramActionDesc)
		progVar := evalProgActionDesc.program
		protoProgVar, err := CreateProtoFromVariable(progVar)
		if err != nil {
			return nil, err
		}
		evalProgramAction := &protostack.EvalProgramAction{ProgramVariable: protoProgVar.GetProgram()}
		return &protostack.Action{
				Type:       protostack.ActionType_PROG_EVALUATION,
				OpCode:     action.OpCode(),
				RealAction: &protostack.Action_EvalProgramAction{EvalProgramAction: evalProgramAction},
			},
			nil
	}
}

func (e *EvalProgramActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		programVar, err := CreateProgramVariableFromProto(reg, protoAction.GetEvalProgramAction().ProgramVariable)
		if err != nil {
			return nil, err
		}
		return &EvalProgramActionDesc{program: programVar}, nil
	}
}

type VariableDeclarationActionDesc struct {
	varNames        []string
	programVariable *ProgramVariable
}

func (a *VariableDeclarationActionDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {

		variableDeclarationActionDesc := action.(*VariableDeclarationActionDesc)
		progVar := variableDeclarationActionDesc.programVariable
		protoProgVar, err := CreateProtoFromVariable(progVar)
		if err != nil {
			return nil, err
		}
		protoVariableDeclaration := &protostack.VariableDeclarationAction{
			VarNames:        variableDeclarationActionDesc.varNames,
			ProgramVariable: protoProgVar.GetProgram(),
		}
		return &protostack.Action{
				Type:   protostack.ActionType_VARIABLE_DECLARATION,
				OpCode: action.OpCode(),
				RealAction: &protostack.Action_VariableDeclarationAction{
					VariableDeclarationAction: protoVariableDeclaration,
				},
			},
			nil
	}
}

func (a *VariableDeclarationActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		programVar, err := CreateProgramVariableFromProto(reg, protoAction.GetVariableDeclarationAction().GetProgramVariable())
		if err != nil {
			return nil, err
		}
		return &VariableDeclarationActionDesc{
				varNames:        protoAction.GetVariableDeclarationAction().GetVarNames(),
				programVariable: programVar,
			},
			nil
	}
}

// VariableDeclarationActionDesc implements Action
var _ Action = (*VariableDeclarationActionDesc)(nil)

func (a *VariableDeclarationActionDesc) OpCode() string {
	return "__hidden__" + "VariableDeclarationProgram"
}

func (a *VariableDeclarationActionDesc) NbArgs() int {
	return 0
}

func (a *VariableDeclarationActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (a *VariableDeclarationActionDesc) Apply(runtimeContext *RuntimeContext) error {
	runtimeContext.EnterNewScope()
	defer func() {
		runtimeContext.LeaveScope()
	}()
	for _, varName := range a.varNames {
		varValue, err := runtimeContext.stack.Pop()
		if err != nil {
			return err
		}
		err = runtimeContext.SetVariableValue(varName, varValue)
		if err != nil {
			return err
		}
	}
	err := runtimeContext.RunAction(&EvalProgramActionDesc{program: a.programVariable})
	if err != nil {
		return err
	}
	return nil
}

func (a *VariableDeclarationActionDesc) Display() string {
	return fmt.Sprintf("-> %s %s", strings.Join(a.varNames, " "), a.programVariable.display())
}

type EvalActionDesc struct{}

// EvalActionDesc implements Action
var _ Action = (*EvalActionDesc)(nil)

func (e *EvalActionDesc) Display() string {
	return "eval"
}

func (e *EvalActionDesc) OpCode() string {
	return "eval"
}

func (e *EvalActionDesc) NbArgs() int {
	return 1
}

func (e *EvalActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (e *EvalActionDesc) Apply(runtimeContext *RuntimeContext) error {

	v1, err := runtimeContext.stack.Pop()
	if err != nil {
		return err
	}
	switch v1.getType() {
	case TYPE_NUMERIC | TYPE_BOOL | TYPE_STR:
		runtimeContext.stack.Push(v1)
	case TYPE_PROGRAM:
		return executeProgram(runtimeContext, v1.(*ProgramVariable))
	case TYPE_ALG_EXPR:
		expression, err := evalAlgExpression(runtimeContext, v1.(*AlgebraicExpressionVariable).rootNode)
		if err != nil {
			return err
		} else {
			runtimeContext.stack.Push(expression)
			return nil
		}
	}

	return nil
}

func evalAlgExpression(runtimeContext *RuntimeContext, algExpreNode AlgebraicExpressionNode) (*NumericVariable, error) {
	numericVariable, _ := algExpreNode.Evaluate(runtimeContext)
	return numericVariable, nil
}

func (e *EvalActionDesc) MarshallFunc() ActionMarshallFunc {
	//TODO implement me
	panic("implement me")
}

func (e *EvalActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	//TODO implement me
	panic("implement me")
}

var evalAct = &EvalActionDesc{}

var StructOpsPackage = ActionPackage{
	staticActions: []Action{
		evalAct,
	},
	dynamicActions: []Action{
		&EvalProgramActionDesc{},
		&IfThenElseActionDesc{},
		&ForNextLoopActionDesc{},
		&StartNextLoopActionDesc{},
		&VariableDeclarationActionDesc{},
		&VariableEvaluationActionDesc{},
		&VariablePutOnStackActionDesc{},
	},
}
