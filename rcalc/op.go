package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

func CheckNoop(elts ...Variable) (bool, error) {
	return true, nil
}

func CheckFirstInt(elts ...Variable) (bool, error) {

	if elts[0].getType() != TYPE_NUMERIC {
		return false, nil
	} else {
		v := elts[0].asNumericVar().value
		if !v.IsInteger() {
			return false, fmt.Errorf("%v is not an integer", v)
		}
	}
	return true, nil
}

func CheckGen(types []Type) CheckTypeFn {
	return func(elts ...Variable) (bool, error) {
		var errors []string

		for idx, varType := range types {

			observedType := elts[idx].getType()
			if varType == TYPE_GENERIC {
				break
			} else if observedType != varType {
				errors = append(errors, fmt.Sprintf("Type error at level %d, expected: %v, found: %v", idx+1, observedType, varType))
			}
		}
		if len(errors) == 0 {
			return true, nil
		} else {
			return false, fmt.Errorf("%s", strings.Join(errors, "\n"))
		}
	}
}

func NewStackOp(opCode string, nbArgs int, fn PureOperationApplyFn) OperationDesc {
	return NewOperationDesc(opCode, nbArgs, CheckNoop, OpToActionFn(fn))
}

func NewStackOpWithtypeCheck(opCode string, nbArgs int, checkFn CheckTypeFn, fn PureOperationApplyFn) OperationDesc {
	return NewOperationDesc(opCode, nbArgs, checkFn, OpToActionFn(fn))
}

func NewRawStackOpWithCheck(opCode string, nbArgs int, checkFn CheckTypeFn, fn ActionApplyFn) ActionDesc {
	return NewActionDesc(opCode, nbArgs, checkFn, fn)
}

// Tooling for Numeric (Decimal) functions

func GetEltAsNumeric(elts []Variable, idx int) decimal.Decimal {
	return elts[idx].asNumericVar().value
}

func CheckAllNumerics(elts ...Variable) (bool, error) {
	for _, e := range elts {
		if e.getType() != TYPE_NUMERIC {
			return false, nil
		}
	}
	return true, nil
}

type A1R1NumericFn func(num1 decimal.Decimal) decimal.Decimal

func A1R1NumericApplyFn(f A1R1NumericFn) PureOperationApplyFn {
	return func(elts ...Variable) []Variable {
		elt := GetEltAsNumeric(elts, 0)
		return []Variable{CreateNumericVariable(f(elt))}
	}
}

func NewA1R1NumericOp(opCode string, decimalFunc A1R1NumericFn) OperationDesc {
	return NewOperationDesc(opCode, 1, CheckAllNumerics, OpToActionFn(A1R1NumericApplyFn(decimalFunc)))
}

type A2R1NumericFn func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal

func A2R1NumericApplyFn(f A2R1NumericFn) PureOperationApplyFn {
	return func(elts ...Variable) []Variable {
		elt1 := GetEltAsNumeric(elts, 1)
		elt2 := GetEltAsNumeric(elts, 0)
		return []Variable{CreateNumericVariable(f(elt1, elt2))}
	}
}

func NewA2R1NumericOp(opCode string, decimalFunc A2R1NumericFn) OperationDesc {
	return NewOperationDesc(opCode, 2, CheckAllNumerics, OpToActionFn(A2R1NumericApplyFn(decimalFunc)))
}

func GetEltAsBoolean(elts []Variable, idx int) bool {
	return elts[idx].asBooleanVar().value
}

func CheckAllBooleans(elts ...Variable) (bool, error) {
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

func A1R1BooleanApplyFn(f A1R1BooleanFn) PureOperationApplyFn {
	return func(elts ...Variable) []Variable {
		elt := GetEltAsBoolean(elts, 0)
		return []Variable{CreateBooleanVariable(f(elt))}
	}
}

func NewA1R1BooleanOp(opCode string, booleanFunc A1R1BooleanFn) OperationDesc {
	return NewOperationDesc(opCode, 1, CheckAllBooleans, OpToActionFn(A1R1BooleanApplyFn(booleanFunc)))
}

type A2R1BooleanFn func(b1 bool, b2 bool) bool

func A2R1BooleanApplyFn(f A2R1BooleanFn) PureOperationApplyFn {
	return func(elts ...Variable) []Variable {
		elt := GetEltAsBoolean(elts, 1)
		elt2 := GetEltAsBoolean(elts, 0)
		return []Variable{CreateBooleanVariable(f(elt, elt2))}
	}
}

func NewA2R1BooleanOp(opCode string, booleanFunc A2R1BooleanFn) OperationDesc {
	return NewOperationDesc(opCode, 2, CheckAllBooleans, OpToActionFn(A2R1BooleanApplyFn(booleanFunc)))
}

type A2NumericR1BooleanFn func(d1 decimal.Decimal, d2 decimal.Decimal) bool

func A2NumericR1BooleanApplyFn(f A2NumericR1BooleanFn) PureOperationApplyFn {
	return func(elts ...Variable) []Variable {
		elt := GetEltAsNumeric(elts, 1)
		elt2 := GetEltAsNumeric(elts, 0)
		return []Variable{CreateBooleanVariable(f(elt, elt2))}
	}
}

func NewA2NumericR1BooleanOp(opCode string, numericToBooleanFunc A2NumericR1BooleanFn) OperationDesc {
	return NewOperationDesc(opCode,
		2,
		CheckAllNumerics,
		OpToActionFn(A2NumericR1BooleanApplyFn(numericToBooleanFunc)))
}

var VersionOp = NewOperationDesc(
	"VERSION",
	0,
	func(elts ...Variable) (bool, error) { return true, nil },
	OpToActionFn(func(elts ...Variable) []Variable {
		return []Variable{CreateNumericVariable(decimal.Zero)}
	}))

var DebugOp = NewOperationDesc(
	"debug",
	1,
	func(elts ...Variable) (bool, error) { return true, nil },
	OpToActionFn(func(elts ...Variable) []Variable {
		fmt.Printf("%v\n", elts[0])
		return []Variable{
			elts[0]}
	}))
