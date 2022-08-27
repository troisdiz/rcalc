package rcalc

import (
	"fmt"
	"strings"
)

type VariablePutOnStackActionDesc struct {
	value Variable
}

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

type VariableEvaluationActionDesc struct {
	varName string
}

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

type StartNextLoopActionDesc struct {
	actions []Action
}

func (a *StartNextLoopActionDesc) NbArgs() int {
	return 2
}

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

func CreateStartNextLoopAction(actions []Action) *StartNextLoopActionDesc {
	return &StartNextLoopActionDesc{actions: actions}
}

type ForNextLoopActionDesc struct {
	varName string
	actions []Action
}

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

type EvalProgramActionDesc struct {
	program *ProgramVariable
}

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
	runtimeContext.EnterNewScope()
	defer func() { runtimeContext.LeaveScope() }()

	for _, action := range e.program.actions {
		err := runtimeContext.RunAction(action)
		if err != nil {
			return err
		}
	}
	return nil
}

type VariableDeclarationActionDesc struct {
	varNames        []string
	programVariable *ProgramVariable
}

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
