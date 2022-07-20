package rcalc

import "fmt"

type Action interface {
	OpCode() string
	NbArgs() int
	CheckTypes(elts ...Variable) (bool, error)
	Apply(system System, stack *Stack) error
}

type ActionCommonDesc struct {
	opCode string
}

func (op *ActionCommonDesc) OpCode() string {
	return op.opCode
}

// ActionDesc implementation of Action interface
type ActionApplyFn func(system System, stack *Stack) error

type ActionDesc struct {
	ActionCommonDesc
	nbArgs        int
	actionApplyFn ActionApplyFn
}

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

func (a *ActionDesc) Apply(system System, stack *Stack) error {
	return a.actionApplyFn(system, stack)
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

func (op *OperationDesc) Apply(system System, stack *Stack) error {
	inputs, err := stack.PopN(op.NbArgs())
	if err != nil {
		return err
	}
	results := op.applyFn(system, inputs...)
	for _, elt := range results {
		stack.Push(elt)
	}
	return nil
}

type PureOperationApplyFn func(elts ...Variable) []Variable

func OpToActionFn(opFn PureOperationApplyFn) OperationApplyFn {
	return func(system System, elts ...Variable) []Variable {
		return opFn(elts...)
	}
}

/* Registry stuff */

type ActionRegistry struct {
	actionDescs map[string]Action
}

func (reg *ActionRegistry) Register(aDesc Action) {
	reg.actionDescs[aDesc.OpCode()] = aDesc
}

type ActionPackage struct {
	actions []Action
}

func (reg *ActionRegistry) RegisterActions(aPackage *ActionPackage) {
	for _, aDesc := range aPackage.actions {
		reg.actionDescs[aDesc.OpCode()] = aDesc
	}
}

func initRegistry() *ActionRegistry {
	reg := ActionRegistry{
		actionDescs: map[string]Action{},
	}
	reg.RegisterActions(&ArithmeticPackage)
	reg.RegisterActions(&TrigonometricPackage)
	reg.RegisterActions(&BooleanLogicPackage)
	reg.RegisterActions(&StackPackage)
	reg.RegisterActions(&MemoryPackage)
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

var Registry = initRegistry()
