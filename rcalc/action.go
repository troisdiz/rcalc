package rcalc

import "fmt"

type Action interface {
	OpCode() string
	NbArgs() int
	CheckTypes(elts ...StackElt) (bool, error)
	Apply(system System, stack *Stack) error
}

type ActionCommonDesc struct {
	opCode string
}

func (op *ActionCommonDesc) OpCode() string {
	return op.opCode
}

// OperationDesc implementation of Action interface
type OperationCommonDesc struct {
	ActionCommonDesc
}

type ActionApplyFn func(system System, elts ...StackElt) []StackElt

type OperationDesc struct {
	OperationCommonDesc
	nbArgs      int
	checkTypeFn CheckTypeFn
	applyFn     ActionApplyFn
}

func NewOperationDesc(opCode string, nbArgs int, checkTypeFn CheckTypeFn, applyFn ActionApplyFn) OperationDesc {
	return OperationDesc{
		OperationCommonDesc: OperationCommonDesc{
			ActionCommonDesc: ActionCommonDesc{
				opCode: opCode,
			},
		},
		nbArgs:      nbArgs,
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

func (op *OperationDesc) CheckTypes(elts ...StackElt) (bool, error) {
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

type ActionRegistry struct {
	actionDescs map[string]*OperationDesc
}

func (reg *ActionRegistry) Register(aDesc *OperationDesc) {
	reg.actionDescs[aDesc.opCode] = aDesc
}

type ActionPackage struct {
	actions []*OperationDesc
}

func (reg *ActionRegistry) RegisterActions(aPackage *ActionPackage) {
	for _, aDesc := range aPackage.actions {
		reg.actionDescs[aDesc.opCode] = aDesc
	}
}

func initRegistry() *ActionRegistry {
	reg := ActionRegistry{
		actionDescs: map[string]*OperationDesc{},
	}
	reg.RegisterActions(&ArithmeticPackage)
	reg.RegisterActions(&TrigonometricPackage)
	reg.RegisterActions(&BooleanLogicPackage)
	reg.RegisterActions(&StackPackage)
	reg.Register(&VersionOp)
	reg.Register(&EXIT_ACTION)
	return &reg
}

func (reg *ActionRegistry) ContainsOpCode(opCode string) bool {
	_, ok := reg.actionDescs[opCode]
	return ok
}

func (reg *ActionRegistry) GetAction(opCode string) *OperationDesc {
	actionDesc, ok := reg.actionDescs[opCode]
	if !ok {
		return nil
	} else {
		return actionDesc
	}
}

var Registry = initRegistry()
