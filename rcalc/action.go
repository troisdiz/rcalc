package rcalc

import "fmt"

type ActionApplyFn func(system System, elts ...StackElt) []StackElt

type Action interface {
	NbArgs() int
	CheckTypes(elts ...StackElt) (bool, error)
	Apply(system System, elts ...StackElt) []StackElt
}

// ActionDesc implementation of Action interface
type ActionDesc struct {
	opCode      string
	nbArgs      int
	checkTypeFn CheckTypeFn
	applyFn     ActionApplyFn
}

func (op *ActionDesc) String() string {
	return fmt.Sprintf("Action(opCode = %s, nbArgs = %d)", op.opCode, op.nbArgs)
}

func (op *ActionDesc) OpCode() string {
	return op.opCode
}

func (op *ActionDesc) NbArgs() int {
	return op.nbArgs
}

func (op *ActionDesc) CheckTypes(elts ...StackElt) (bool, error) {
	return op.checkTypeFn(elts...)
}

func (op *ActionDesc) Apply(system System, elts ...StackElt) []StackElt {
	return op.applyFn(system, elts...)
}

type ActionRegistry struct {
	actionDescs map[string]*ActionDesc
}

func (reg *ActionRegistry) Register(aDesc *ActionDesc) {
	reg.actionDescs[aDesc.opCode] = aDesc
}

type ActionPackage struct {
	actions []*ActionDesc
}

func (reg *ActionRegistry) RegisterActions(aPackage *ActionPackage) {
	for _, aDesc := range aPackage.actions {
		reg.actionDescs[aDesc.opCode] = aDesc
	}
}

func initRegistry() *ActionRegistry {
	reg := ActionRegistry{
		actionDescs: map[string]*ActionDesc{},
	}
	reg.RegisterActions(&ArithmeticPackage)
	reg.RegisterActions(&TrigonometricPackage)
	reg.RegisterActions(&BooleanLogicPackage)
	reg.Register(&VersionOp)
	reg.Register(&EXIT_ACTION)
	return &reg
}

func (reg *ActionRegistry) ContainsOpCode(opCode string) bool {
	_, ok := reg.actionDescs[opCode]
	return ok
}

func (reg *ActionRegistry) GetAction(opCode string) *ActionDesc {
	actionDesc, ok := reg.actionDescs[opCode]
	if !ok {
		return nil
	} else {
		return actionDesc
	}
}

var Registry = initRegistry()
