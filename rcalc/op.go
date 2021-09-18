package rcalc

import "github.com/shopspring/decimal"

type CheckTypeFn func(elts ...StackElt) (bool, error)
type OpApplyFn func(elts ...StackElt) []StackElt

func OpToActionFn(opFn OpApplyFn) ActionApplyFn {
	return func(system System, elts ...StackElt) []StackElt {
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

type NumOp2Args1Result func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal

func DecimalFuncToOpApplyFn(f NumOp2Args1Result) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt1 := GetEltAsNumeric(elts, 1)
		elt2 := GetEltAsNumeric(elts, 0)
		return []StackElt{ CreateNumericStackElt(f(elt1, elt2)) }
	}
}

func NewTwoArgsSingleResultNumOp(opCode string, decimalFunc NumOp2Args1Result) ActionDesc {
	return ActionDesc{
		opCode: opCode,
		nbArgs: 2,
		checkTypeFn: Check2Numerics,
		applyFn: OpToActionFn(DecimalFuncToOpApplyFn(decimalFunc)),
	}

}

var addOp = NewTwoArgsSingleResultNumOp("+", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Add(num2)
})

var subOp = NewTwoArgsSingleResultNumOp("-", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Sub(num2)
})

var mulOp = NewTwoArgsSingleResultNumOp("*", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Mul(num2)
})

var divOp = NewTwoArgsSingleResultNumOp("/", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Div(num2)
})

var powOp = NewTwoArgsSingleResultNumOp("^", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Pow(num2)
})

var VersionOp = ActionDesc{
	opCode:      "VERSION",
	nbArgs:      0,
	checkTypeFn: func(elts ...StackElt) (bool, error) { return true, nil },
	applyFn:     OpToActionFn(func(elts ...StackElt) []StackElt {
		return []StackElt{ CreateNumericStackElt(decimal.Zero) }
	}),
}
