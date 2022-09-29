package rcalc

import (
	"fmt"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

type ActionMarshallFunc func(reg *ActionRegistry, action Action) (*protostack.Action, error)
type ActionUnMarshallFunc func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error)

type Action interface {
	OpCode() string
	NbArgs() int
	CheckTypes(elts ...Variable) (bool, error)
	Apply(runtimeContext *RuntimeContext) error
	Display() string
	MarshallFunc() ActionMarshallFunc
	UnMarshallFunc() ActionUnMarshallFunc
}

type ActionCommonDesc struct {
	opCode string
}

func (op *ActionCommonDesc) OpCode() string {
	return op.opCode
}
func ActionCommonDescMarshalFunction(reg *ActionRegistry, action Action) (*protostack.Action, error) {
	return &protostack.Action{
		Type:   protostack.ActionType_OPERATION,
		OpCode: action.OpCode(),
	}, nil
}

func (op *ActionCommonDesc) MarshallFunc() ActionMarshallFunc {
	return ActionCommonDescMarshalFunction
}

func ActionCommonDescUnMarshallFunc(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
	protoOpCode := protoAction.OpCode
	if reg.ContainsOpCode(protoOpCode) {
		return reg.GetAction(protoOpCode), nil
	} else {
		return nil, fmt.Errorf("cannot find action with opcode %s", protoOpCode)
	}
}

func (op *ActionCommonDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return ActionCommonDescUnMarshallFunc
}

// ActionDesc implementation of Action interface
type ActionApplyFn func(system System, stack *Stack) error

type ActionDesc struct {
	ActionCommonDesc
	nbArgs        int
	actionApplyFn ActionApplyFn
}

func (a *ActionDesc) MarshallFunc() ActionMarshallFunc {
	//TODO implement me
	panic("implement me ActionDesc MarshallFunc")
}

func (a *ActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	//TODO implement me
	panic("implement me ActionDesc UnMarshallFunc")
}

var _ Action = (*ActionDesc)(nil)

func NewActionDesc(opCode string, nbArgs int, checkTypeFn CheckTypeFn, applyFn ActionApplyFn) ActionDesc {
	return ActionDesc{
		ActionCommonDesc: ActionCommonDesc{
			opCode: opCode,
		},
		nbArgs:        nbArgs,
		actionApplyFn: applyFn,
	}
}
func (a *ActionDesc) NbArgs() int {
	return a.nbArgs
}

func (a *ActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return CheckNoop(elts...)
}

func (a *ActionDesc) Apply(runtimeContext *RuntimeContext) error {
	return a.actionApplyFn(runtimeContext.system, runtimeContext.stack)
}

func (a *ActionDesc) Display() string {
	return a.OpCode()
}

// OperationDesc implementation of Action interface
type OperationCommonDesc struct {
	ActionCommonDesc
	nbArgs int
}

type CheckTypeFn func(elts ...Variable) (bool, error)
type OperationApplyFn func(system System, elts ...Variable) []Variable

type OperationDesc struct {
	OperationCommonDesc
	checkTypeFn CheckTypeFn
	applyFn     OperationApplyFn
}

var _ Action = (*OperationDesc)(nil)

func NewOperationDesc(opCode string, nbArgs int, checkTypeFn CheckTypeFn, applyFn OperationApplyFn) OperationDesc {
	return OperationDesc{
		OperationCommonDesc: OperationCommonDesc{
			ActionCommonDesc: ActionCommonDesc{
				opCode: opCode,
			},
			nbArgs: nbArgs,
		},
		checkTypeFn: checkTypeFn,
		applyFn:     applyFn,
	}
}

func (op *OperationDesc) String() string {
	return fmt.Sprintf("Action(opCode = %s, nbArgs = %d)", op.opCode, op.nbArgs)
}

func (op *OperationDesc) NbArgs() int {
	return op.nbArgs
}

func (op *OperationDesc) CheckTypes(elts ...Variable) (bool, error) {
	return op.checkTypeFn(elts...)
}

func (op *OperationDesc) Apply(runtimeContext *RuntimeContext) error {
	inputs, err := runtimeContext.stack.PopN(op.NbArgs())
	if err != nil {
		return err
	}
	results := op.applyFn(runtimeContext.system, inputs...)
	for _, elt := range results {
		runtimeContext.stack.Push(elt)
	}
	return nil
}

func (op *OperationDesc) Display() string {
	return op.opCode
}

/*
func (op *OperationDesc) MarshallFunc() ActionMarshallFunc {
	return func(reg *ActionRegistry, action Action) (*protostack.Action, error) {
		op := action.(*OperationDesc)
		return &protostack.Action{
			Type:   protostack.ActionType_OPERATION,
			OpCode: op.OpCode(),
		}, nil
	}
}

func (op *OperationDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
		protoOpCode := protoAction.OpCode
		if reg.ContainsOpCode(protoOpCode) {
			return reg.GetAction(protoOpCode), nil
		} else {
			return nil, fmt.Errorf("cannot find action with opcode %s", protoOpCode)
		}
	}
}*/

type PureOperationApplyFn func(elts ...Variable) []Variable

func OpToActionFn(opFn PureOperationApplyFn) OperationApplyFn {
	return func(system System, elts ...Variable) []Variable {
		return opFn(elts...)
	}
}

/* Registry stuff */

type ActionRegistry struct {
	actionDescs    map[string]Action
	dynamicActions map[string]struct {
		marshalFunc   ActionMarshallFunc
		unMarshalFunc ActionUnMarshallFunc
	}
	// marshallFunctions   map[string]ActionMarshallFunc
	// unMarshallFunctions map[string]ActionUnMarshallFunc
}

func (reg *ActionRegistry) Register(aDesc Action) {
	reg.actionDescs[aDesc.OpCode()] = aDesc
}

type ActionPackage struct {
	staticActions  []Action
	dynamicActions []Action
}

func (ap *ActionPackage) AddStatic(action Action) {
	ap.staticActions = append(ap.staticActions, action)
}

func (ap *ActionPackage) AddDynamic(action Action) {
	ap.dynamicActions = append(ap.dynamicActions, action)
}

func (reg *ActionRegistry) RegisterActions(aPackage *ActionPackage) {
	for _, aDesc := range aPackage.staticActions {
		reg.actionDescs[aDesc.OpCode()] = aDesc
	}
	for _, dynAction := range aPackage.dynamicActions {
		// reg.marshallFunctions[dynAction.OpCode()] = dynAction.MarshallFunc()
		// reg.unMarshallFunctions[dynAction.OpCode()] = dynAction.UnMarshallFunc()

		reg.dynamicActions[dynAction.OpCode()] = struct {
			marshalFunc   ActionMarshallFunc
			unMarshalFunc ActionUnMarshallFunc
		}{
			marshalFunc:   dynAction.MarshallFunc(),
			unMarshalFunc: dynAction.UnMarshallFunc(),
		}
	}
}

func initRegistry() *ActionRegistry {
	reg := ActionRegistry{
		actionDescs: map[string]Action{},
		dynamicActions: map[string]struct {
			marshalFunc   ActionMarshallFunc
			unMarshalFunc ActionUnMarshallFunc
		}{},
		// marshallFunctions:   map[string]ActionMarshallFunc{},
		// unMarshallFunctions: map[string]ActionUnMarshallFunc{},
	}
	reg.RegisterActions(&ArithmeticPackage)
	reg.RegisterActions(&TrigonometricPackage)
	reg.RegisterActions(&BooleanLogicPackage)
	reg.RegisterActions(&StackPackage)
	reg.RegisterActions(&MemoryPackage)
	reg.RegisterActions(&StructOpsPackage)
	reg.Register(&VersionOp)
	reg.Register(&EXIT_ACTION)
	return &reg
}

func (reg *ActionRegistry) ContainsOpCode(opCode string) bool {
	_, ok := reg.actionDescs[opCode]
	return ok
}

func (reg *ActionRegistry) GetAction(opCode string) Action {
	actionDesc, ok := reg.actionDescs[opCode]
	if !ok {
		return nil
	} else {
		return actionDesc
	}
}

func (reg *ActionRegistry) GetDynamicActionMarshallFunc(opCode string) ActionMarshallFunc {
	if dynAction, ok := reg.dynamicActions[opCode]; ok {
		return dynAction.marshalFunc
	} else {
		return nil
	}
}

func (reg *ActionRegistry) GetDynamicActionUnMarshallFunc(opCode string) ActionUnMarshallFunc {
	if dynAction, ok := reg.dynamicActions[opCode]; ok {
		return dynAction.unMarshalFunc
	} else {
		return nil
	}
}

func (reg *ActionRegistry) CreateActionFromProto(protoAction *protostack.Action) (Action, error) {
	if mFuncs, ok := reg.dynamicActions[protoAction.OpCode]; ok {
		action, err := mFuncs.unMarshalFunc(reg, protoAction)
		if err != nil {
			return nil, err
		}
		return action, nil
	}
	return nil, fmt.Errorf("no unMarshallFunction found for type %d / OpCode %s", protoAction.Type, protoAction.OpCode)
}

func (reg *ActionRegistry) GetDynamicActionOpCodes() []string {
	result := make([]string, len(reg.dynamicActions))
	for k := range reg.dynamicActions {
		result = append(result, k)
	}
	return result
}

var Registry = initRegistry()
