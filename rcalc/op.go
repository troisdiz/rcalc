package rcalc

type CheckTypeFn func(elts ...StackElt) (bool, error)
type OpApplyFn func(elts ...StackElt) StackElt

func OpToActionFn(opFn OpApplyFn) ActionApplyFn {
	return func(system System, elts ...StackElt) StackElt {
		return opFn(elts...)
	}
}

func GetEltAsInt(elts []StackElt, idx int) int {
	return elts[idx].asIntElt().value
}

func Check2Ints(elts ...StackElt) (bool, error) {
	for _, e := range elts {
		if  e.getType() != TYPE_INT {
			return false, nil
		}
	}
	return true, nil
}

var ADD_OP = ActionDesc{
	opCode:      "+",
	nbArgs:      2,
	checkTypeFn: Check2Ints,
	applyFn:     OpToActionFn(func(elt ...StackElt) StackElt {
		elt1 := GetEltAsInt(elt, 0)
		elt2 := GetEltAsInt(elt, 1)
		return CreateIntStackElt(elt1 + elt2)
	}),
}

var MUL_OP = ActionDesc{
	opCode:      "*",
	nbArgs:      2,
	checkTypeFn: Check2Ints,
	applyFn:     OpToActionFn(func(elt ...StackElt) StackElt {
		elt1 := GetEltAsInt(elt, 0)
		elt2 := GetEltAsInt(elt, 1)
		return CreateIntStackElt(elt1 * elt2)
	}),
}

var VERSION_OP = ActionDesc{
	opCode:      "VERSION",
	nbArgs:      0,
	checkTypeFn: func(elts ...StackElt) (bool, error) { return true, nil },
	applyFn:     OpToActionFn(func(elts ...StackElt) StackElt { return CreateIntStackElt(0) }),
}

