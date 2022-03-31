package rcalc

import "github.com/shopspring/decimal"

// Arithmetic package

var addOp = NewA2R1NumericOp("+", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Add(num2)
})

var subOp = NewA2R1NumericOp("-", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Sub(num2)
})

var mulOp = NewA2R1NumericOp("*", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Mul(num2)
})

var divOp = NewA2R1NumericOp("/", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Div(num2)
})

var powOp = NewA2R1NumericOp("^", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Pow(num2)
})

var ArithmeticPackage = ActionPackage{
	[]*ActionDesc{&addOp, &subOp, &mulOp, &divOp, &powOp},
}

// Trigonometry package

var sinOp = NewA1R1NumericOp("sin", func(num decimal.Decimal) decimal.Decimal {
	return num.Sin()
})

var arcSinOp = NewA1R1NumericOp("asin", func(num decimal.Decimal) decimal.Decimal {
	return num.Div(decimal.NewFromInt(1).Sub(num.Pow(decimal.NewFromInt(2))).Pow(decimal.New(5, -1))).Atan()
})

var cosOp = NewA1R1NumericOp("cos", func(num decimal.Decimal) decimal.Decimal {
	return num.Cos()
})

var arcCosOp = NewA1R1NumericOp("acos", func(num decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(1).Sub(num.Pow(decimal.NewFromInt(2))).Pow(decimal.New(5, -1)).Div(num).Atan()
})

var tanOp = NewA1R1NumericOp("sin", func(num decimal.Decimal) decimal.Decimal {
	return num.Tan()
})

var arcTanOp = NewA1R1NumericOp("atan", func(num decimal.Decimal) decimal.Decimal {
	return num.Atan()
})

var TrigonometricPackage = ActionPackage{
	[]*ActionDesc{
		&sinOp, &cosOp, &tanOp,
		&arcSinOp, &arcCosOp, &arcTanOp,
	},
}

// Logic Package
var negOp = NewA1R1BooleanOp("neg", func(b bool) bool {
	return !b
})

var andOp = NewA2R1BooleanOp("and", func(b bool, b2 bool) bool {
	return b && b2
})

var orOp = NewA2R1BooleanOp("or", func(b bool, b2 bool) bool {
	return b || b2
})

var xorOp = NewA2R1BooleanOp("xor", func(b bool, b2 bool) bool {
	return b != b2
})

var xandOp = NewA2R1BooleanOp("xand", func(b bool, b2 bool) bool {
	return b == b2
})

var BooleanLogicPackage = ActionPackage{
	[]*ActionDesc{
		&negOp, &andOp, &orOp, &xorOp, &xandOp,
	},
}

// Stack package

var dupOp = NewStackOp("dup", 1, func(elts ...StackElt) []StackElt {
	return []StackElt{elts[0], elts[0]}
})

var dupNOp = NewStackOpWithtypeCheck("dupn", 1, CheckFirstInt, func(elts ...StackElt) []StackElt {
	var result []StackElt = make([]StackElt, elts[0].asNumericElt().value.IntPart())
	return result
})

var StackPackage = ActionPackage{
	[]*ActionDesc{
		&dupOp,
	},
}
