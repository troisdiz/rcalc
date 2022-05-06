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
	[]Action{&addOp, &subOp, &mulOp, &divOp, &powOp},
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
	[]Action{
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
	[]Action{
		&negOp, &andOp, &orOp, &xorOp, &xandOp,
	},
}

// Stack package

var dupOp = NewStackOp("dup", 1, func(elts ...StackElt) []StackElt {
	return []StackElt{elts[0], elts[0]}
})

var dup2Op = NewStackOp("dup2", 2, func(elts ...StackElt) []StackElt {
	return []StackElt{elts[1], elts[0], elts[1], elts[0]}
})

var dropOp = NewStackOp("drop", 1, func(elts ...StackElt) []StackElt {
	return []StackElt{}
})

var drop2Op = NewStackOp("drop2", 2, func(elts ...StackElt) []StackElt {
	return []StackElt{}
})

var swapOp = NewStackOp("swap", 2, func(elts ...StackElt) []StackElt {
	return []StackElt{elts[0], elts[1]}
})

// rot, roll, pick

var dupNOp = NewRawStackOpWithCheck("dupn", 1, CheckFirstInt, func(system System, stack *Stack) error {
	n, err := stack.Pop()
	if err != nil {
		return err
	}
	stackElts, err := stack.PopN(int(n.asNumericElt().value.IntPart()))
	if err != nil {
		return err
	}
	stack.PushN(stackElts)
	return nil
})

var depthAct = NewRawStackOpWithCheck("depth", 0, CheckNoop, func(system System, stack *Stack) error {
	stack.Push(CreateNumericStackEltFromInt(stack.Size()))
	return nil
})

var StackPackage = ActionPackage{
	[]Action{
		&dupOp,
		&dup2Op,
		&dropOp,
		&drop2Op,
		&swapOp,
		&dupNOp,
		&depthAct,
	},
}
