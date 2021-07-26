package rcalc

type CheckTypeFn func(elts ...StackElt) (bool, error)
type OpApplyFn func(elts ...StackElt) StackElt

func OpToActionFn(opFn OpApplyFn) ActionApplyFn {
	return func(system System, elts ...StackElt) StackElt {
		return opFn(elts...)
	}
}

var ADD_OP = ActionDesc{
	opCode:      "+",
	nbArgs:      2,
	checkTypeFn: addCheckTypes,
	applyFn:     OpToActionFn(func(elt ...StackElt) StackElt {
		elt1 := elt[0].asIntElt().value
		elt2 := elt[1].asIntElt().value
		return CreateInStackElt(elt1 + elt2)
	}),
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

var VERSION_OP = ActionDesc{
	opCode:      "VERSION",
	nbArgs:      0,
	checkTypeFn: func(elts ...StackElt) (bool, error) { return true, nil },
	applyFn:     OpToActionFn(func(elts ...StackElt) StackElt { return CreateInStackElt(0) }),
}

