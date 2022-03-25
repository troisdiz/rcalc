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
		if e.getType() != TYPE_NUMERIC {
			return false, nil
		}
	}
	return true, nil
}

type NumOp1Arg1Result func(num1 decimal.Decimal) decimal.Decimal

func Decimal1FuncToOp1ApplyFn(f NumOp1Arg1Result) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt := GetEltAsNumeric(elts, 0)
		return []StackElt{CreateNumericStackElt(f(elt))}
	}
}

func NewOneArgSingleResultNumOp(opCode string, decimalFunc NumOp1Arg1Result) ActionDesc {
	return ActionDesc{
		opCode:      opCode,
		nbArgs:      1,
		checkTypeFn: Check2Numerics,
		applyFn:     OpToActionFn(Decimal1FuncToOp1ApplyFn(decimalFunc)),
	}
}

type NumOp2Args1Result func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal

func Decimal2FuncToOp2ApplyFn(f NumOp2Args1Result) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt1 := GetEltAsNumeric(elts, 1)
		elt2 := GetEltAsNumeric(elts, 0)
		return []StackElt{CreateNumericStackElt(f(elt1, elt2))}
	}
}

func NewTwoArgsSingleResultNumOp(opCode string, decimalFunc NumOp2Args1Result) ActionDesc {
	return ActionDesc{
		opCode:      opCode,
		nbArgs:      2,
		checkTypeFn: Check2Numerics,
		applyFn:     OpToActionFn(Decimal2FuncToOp2ApplyFn(decimalFunc)),
	}

}

func GetEltAsBoolean(elts []StackElt, idx int) bool {
	return elts[idx].asBooleanElt().value
}

func CheckAllBooleans(elts ...StackElt) (bool, error) {
	for _, e := range elts {
		if e.getType() != TYPE_BOOL {
			return false, nil
		}
	}
	return true, nil
}

type BooleanOp1Arg1Result func(num1 bool) bool

func BooleanFuncToOp1Arg1ResultApplyFn(f BooleanOp1Arg1Result) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt := GetEltAsBoolean(elts, 0)
		return []StackElt{CreateBooleanStackElt(f(elt))}
	}
}

func New1Arg1ResultBooleanOp(opCode string, booleanFunc BooleanOp1Arg1Result) ActionDesc {
	return ActionDesc{
		opCode:      opCode,
		nbArgs:      1,
		checkTypeFn: CheckAllBooleans,
		applyFn:     OpToActionFn(BooleanFuncToOp1Arg1ResultApplyFn(booleanFunc)),
	}
}

var VersionOp = ActionDesc{
	opCode:      "VERSION",
	nbArgs:      0,
	checkTypeFn: func(elts ...StackElt) (bool, error) { return true, nil },
	applyFn: OpToActionFn(func(elts ...StackElt) []StackElt {
		return []StackElt{CreateNumericStackElt(decimal.Zero)}
	}),
}
