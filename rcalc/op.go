package rcalc

import "github.com/shopspring/decimal"

type CheckTypeFn func(elts ...StackElt) (bool, error)
type OpApplyFn func(elts ...StackElt) StackElt

func OpToActionFn(opFn OpApplyFn) ActionApplyFn {
	return func(system System, elts ...StackElt) StackElt {
		return opFn(elts...)
	}
}

func GetEltAsNumeric(elts []StackElt, idx int) decimal.Decimal {
	return elts[idx].asNumericElt().value
}

func Check2Numerics(elts ...StackElt) (bool, error) {
	for _, e := range elts {
		if  e.getType() != TYPE_NUMERIC {
			return false, nil
		}
	}
	return true, nil
}

var ADD_OP = ActionDesc{
	opCode:      "+",
	nbArgs:      2,
	checkTypeFn: Check2Numerics,
	applyFn:     OpToActionFn(func(elt ...StackElt) StackElt {
		elt1 := GetEltAsNumeric(elt, 0)
		elt2 := GetEltAsNumeric(elt, 1)
		return CreateNumericStackElt(elt1.Add(elt2))
	}),
}

var MUL_OP = ActionDesc{
	opCode:      "*",
	nbArgs:      2,
	checkTypeFn: Check2Numerics,
	applyFn:     OpToActionFn(func(elt ...StackElt) StackElt {
		elt1 := GetEltAsNumeric(elt, 0)
		elt2 := GetEltAsNumeric(elt, 1)
		return CreateNumericStackElt(elt1.Mul(elt2))
	}),
}

var VERSION_OP = ActionDesc{
	opCode:      "VERSION",
	nbArgs:      0,
	checkTypeFn: func(elts ...StackElt) (bool, error) { return true, nil },
	applyFn:     OpToActionFn(func(elts ...StackElt) StackElt { return CreateNumericStackElt(decimal.Zero) }),
}
