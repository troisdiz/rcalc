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
var negOp = New1A1R1BooleanOp("neg", func(b bool) bool {
	return !b
})

var BooleanLogicPackage = ActionPackage{
	[]*ActionDesc{
		&negOp,
	},
}
