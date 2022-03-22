package rcalc

import "github.com/shopspring/decimal"

// Arithmetic package

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

var ArithmeticPackage = ActionPackage{
	[]*ActionDesc{&addOp, &subOp, &mulOp, &divOp, &powOp},
}

// Trigonometry package

var sinOp = NewOneArgSingleResultNumOp("sin", func(num decimal.Decimal) decimal.Decimal {
	return num.Sin()
})

var arcSinOp = NewOneArgSingleResultNumOp("asin", func(num decimal.Decimal) decimal.Decimal {
	return num.Div(decimal.NewFromInt(1).Sub(num.Pow(decimal.NewFromInt(2))).Pow(decimal.New(5, -1))).Atan()
})

var cosOp = NewOneArgSingleResultNumOp("cos", func(num decimal.Decimal) decimal.Decimal {
	return num.Cos()
})

var arcCosOp = NewOneArgSingleResultNumOp("acos", func(num decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(1).Sub(num.Pow(decimal.NewFromInt(2))).Pow(decimal.New(5, -1)).Div(num).Atan()
})

var tanOp = NewOneArgSingleResultNumOp("sin", func(num decimal.Decimal) decimal.Decimal {
	return num.Tan()
})

var arcTanOp = NewOneArgSingleResultNumOp("atan", func(num decimal.Decimal) decimal.Decimal {
	return num.Atan()
})

var TrigonometricPackage = ActionPackage{
	[]*ActionDesc{
		&sinOp, &cosOp, &tanOp,
		&arcSinOp, &arcCosOp, &arcTanOp,
	},
}

// Logic Package
