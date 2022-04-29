package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func OpToActionFn(opFn OpApplyFn) ActionApplyFn {
	return func(system System, elts ...StackElt) []StackElt {
		return opFn(elts...)
	}
}

func CheckNoop(elts ...StackElt) (bool, error) {
	return true, nil
}

func CheckFirstInt(elts ...StackElt) (bool, error) {

	if elts[0].getType() != TYPE_NUMERIC {
		return false, nil
	} else {
		v := elts[0].asNumericElt().value
		if !v.IsInteger() {
			return false, fmt.Errorf("%v is not an integer", v)
		}
	}
	return true, nil
}

func NewStackOp(opCode string, nbArgs int, fn OpApplyFn) OperationDesc {
	return NewOperationDesc(opCode, nbArgs, CheckNoop, OpToActionFn(fn))
}

func NewStackOpWithtypeCheck(opCode string, nbArgs int, checkFn CheckTypeFn, fn OpApplyFn) OperationDesc {
	return NewOperationDesc(opCode, nbArgs, checkFn, OpToActionFn(fn))
}

// Tooling for Numeric (Decimal) functions

func GetEltAsNumeric(elts []StackElt, idx int) decimal.Decimal {
	return elts[idx].asNumericElt().value
}

func CheckAllNumerics(elts ...StackElt) (bool, error) {
	for _, e := range elts {
		if e.getType() != TYPE_NUMERIC {
			return false, nil
		}
	}
	return true, nil
}

type A1R1NumericFn func(num1 decimal.Decimal) decimal.Decimal

func A1R1NumericApplyFn(f A1R1NumericFn) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt := GetEltAsNumeric(elts, 0)
		return []StackElt{CreateNumericStackElt(f(elt))}
	}
}

func NewA1R1NumericOp(opCode string, decimalFunc A1R1NumericFn) OperationDesc {
	return NewOperationDesc(opCode, 1, CheckAllNumerics, OpToActionFn(A1R1NumericApplyFn(decimalFunc)))
}

type A2R1NumericFn func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal

func A2R1NumericApplyFn(f A2R1NumericFn) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt1 := GetEltAsNumeric(elts, 1)
		elt2 := GetEltAsNumeric(elts, 0)
		return []StackElt{CreateNumericStackElt(f(elt1, elt2))}
	}
}

func NewA2R1NumericOp(opCode string, decimalFunc A2R1NumericFn) OperationDesc {
	return NewOperationDesc(opCode, 2, CheckAllNumerics, OpToActionFn(A2R1NumericApplyFn(decimalFunc)))
}

func GetEltAsBoolean(elts []StackElt, idx int) bool {
	return elts[idx].asBooleanElt().value
}

func CheckAllBooleans(elts ...StackElt) (bool, error) {
	fmt.Printf("CheckAllBooleans %v\n", elts)
	for _, e := range elts {
		if e.getType() != TYPE_BOOL {
			return false, nil
		}
	}
	return true, nil
}

// Tooling for boolean functions

type A1R1BooleanFn func(num1 bool) bool

func A1R1BooleanApplyFn(f A1R1BooleanFn) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt := GetEltAsBoolean(elts, 0)
		return []StackElt{CreateBooleanStackElt(f(elt))}
	}
}

func NewA1R1BooleanOp(opCode string, booleanFunc A1R1BooleanFn) OperationDesc {
	return NewOperationDesc(opCode, 1, CheckAllBooleans, OpToActionFn(A1R1BooleanApplyFn(booleanFunc)))
}

type A2R1BooleanFn func(b1 bool, b2 bool) bool

func A2R1BooleanApplyFn(f A2R1BooleanFn) OpApplyFn {
	return func(elts ...StackElt) []StackElt {
		elt := GetEltAsBoolean(elts, 1)
		elt2 := GetEltAsBoolean(elts, 0)
		return []StackElt{CreateBooleanStackElt(f(elt, elt2))}
	}
}

func NewA2R1BooleanOp(opCode string, booleanFunc A2R1BooleanFn) OperationDesc {
	return NewOperationDesc(opCode, 2, CheckAllBooleans, OpToActionFn(A2R1BooleanApplyFn(booleanFunc)))
}

var VersionOp = NewOperationDesc(
	"VERSION",
	0,
	func(elts ...StackElt) (bool, error) { return true, nil },
	OpToActionFn(func(elts ...StackElt) []StackElt {
		return []StackElt{CreateNumericStackElt(decimal.Zero)}
	}))
