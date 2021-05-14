package rcalc

type CheckTypeFn func(elts ...StackElt) (bool, error)
type ApplyFn func(elts ...StackElt) StackElt

type Operation interface {
	NbArgs() int
	CheckTypes(elts ...StackElt) (bool, error)
	Apply(elts ...*StackElt) *StackElt
}

type OperationDesc struct {
	opCode      string
	nbArgs      int
	checkTypeFn CheckTypeFn
	applyFn     ApplyFn
}

func (op *OperationDesc) OpCode() string {
	return op.opCode
}

func (op *OperationDesc) NbArgs() int {
	return op.nbArgs
}

func (op *OperationDesc) CheckTypes(elts ...StackElt) (bool, error)  {
	return op.checkTypeFn(elts...)
}

func (op *OperationDesc) Apply(elts ...StackElt) StackElt  {
	return op.applyFn(elts...)
}

/* Operation registry */
type OperationRegistry struct {
	ops map[string]*OperationDesc
}

func (reg *OperationRegistry) Register(op *OperationDesc)  {
	reg.ops[op.opCode] = op
}

/* Operations */
var ADD_OP = OperationDesc{
	opCode:      "+",
	nbArgs:      2,
	checkTypeFn: addCheckTypes,
	applyFn:     addApply,
}

func addCheckTypes(elts ...StackElt) (bool, error)  {
	for _, e := range elts {
		if  e.getType() != TYPE_INT {
			return false, nil
		}
	}
	return true, nil
}

func addApply(elt ...StackElt) StackElt  {
	elt1 := elt[0].asIntElt().value
	elt2 := elt[1].asIntElt().value
	return CreateInStackElt(elt1+elt2)
}

var VERSION_OP = OperationDesc{
	opCode:      "VERSION",
	nbArgs:      0,
	checkTypeFn: func(elts ...StackElt) (bool, error) { return true, nil },
	applyFn:     func(elts ...StackElt) StackElt { return CreateInStackElt(0) },
}

func initRegistry() *OperationRegistry {
	reg := OperationRegistry{
		ops: map[string]*OperationDesc{},
	}
	reg.Register(&ADD_OP)
	reg.Register(&VERSION_OP)
	return &reg
}

var Registry = initRegistry()
