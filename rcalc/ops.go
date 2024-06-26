package rcalc

import (
	"github.com/shopspring/decimal"
	"gonum.org/v1/gonum/stat/combin"
)

// Arithmetic package

var addOp = NewExpandedA2R1NumericOp("+", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Add(num2)
})

var subOp = NewExpandedA2R1NumericOp("-", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num2.Sub(num1)
})

var mulOp = NewExpandedA2R1NumericOp("*", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num1.Mul(num2)
})

var divOp = NewExpandedA2R1NumericOp("/", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num2.Div(num1)
})

var powOp = NewExpandedA2R1NumericOp("^", func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
	return num2.Pow(num1)
})

var ArithmeticPackage = ActionPackage{
	staticActions: []Action{&addOp, &subOp, &mulOp, &divOp, &powOp},
}

// Trigonometry package

var sinOp = NewA1R1NumericOp("sin", func(num decimal.Decimal) decimal.Decimal {
	return num.Sin()
})

var sinAlgDesc = AlgebraicFunctionDesc{
	name:      "sin",
	argsCount: 1,
	fn: func(args ...decimal.Decimal) decimal.Decimal {
		return args[0].Sin()
	},
}

var arcSinOp = NewA1R1NumericOp("asin", func(num decimal.Decimal) decimal.Decimal {
	return num.Div(decimal.NewFromInt(1).Sub(num.Pow(decimal.NewFromInt(2))).Pow(decimal.New(5, -1))).Atan()
})

var cosOp = NewA1R1NumericOp("cos", func(num decimal.Decimal) decimal.Decimal {
	return num.Cos()
})

var cosAlgDesc = AlgebraicFunctionDesc{
	name:      "cos",
	argsCount: 1,
	fn: func(args ...decimal.Decimal) decimal.Decimal {
		return args[0].Cos()
	},
}

var arcCosOp = NewA1R1NumericOp("acos", func(num decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(1).Sub(num.Pow(decimal.NewFromInt(2))).Pow(decimal.New(5, -1)).Div(num).Atan()
})

var tanOp = NewA1R1NumericOp("sin", func(num decimal.Decimal) decimal.Decimal {
	return num.Tan()
})

var tanAlgDesc = AlgebraicFunctionDesc{
	name:      "tan",
	argsCount: 1,
	fn: func(args ...decimal.Decimal) decimal.Decimal {
		return args[0].Tan()
	},
}

var arcTanOp = NewA1R1NumericOp("atan", func(num decimal.Decimal) decimal.Decimal {
	return num.Atan()
})

var TrigonometricPackage = ActionPackage{
	staticActions: []Action{
		&sinOp, &cosOp, &tanOp,
		&arcSinOp, &arcCosOp, &arcTanOp,
	},
	algrebraicFunctions: []AlgebraicFunctionDesc{
		sinAlgDesc, cosAlgDesc, tanAlgDesc,
	},
}

// Logic Package
var eqNumOp = NewA2NumericR1BooleanOp("==", func(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d1.Equal(d2)
})

var letNumOp = NewA2NumericR1BooleanOp("<=", func(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d2.LessThanOrEqual(d1)
})

var ltNumOp = NewA2NumericR1BooleanOp("<", func(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d2.LessThan(d1)
})

var getNumOp = NewA2NumericR1BooleanOp(">=", func(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d2.GreaterThanOrEqual(d1)
})

var gtNumOp = NewA2NumericR1BooleanOp("<", func(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d2.GreaterThan(d1)
})

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
	staticActions: []Action{
		&eqNumOp, &ltNumOp, &letNumOp, &gtNumOp, &getNumOp, &negOp, &andOp, &orOp, &xorOp, &xandOp,
	},
}

var combOp = NewA2R1NumericOp("comb", func(p decimal.Decimal, n decimal.Decimal) decimal.Decimal {
	if !n.IsInteger() || !p.IsInteger() {
		return decimal.Zero
	}
	nInt := n.IntPart()
	pInt := p.IntPart()
	return decimal.NewFromInt(int64(combin.Binomial(int(nInt), int(pInt))))
})

var permOp = NewA2R1NumericOp("perm", func(p decimal.Decimal, n decimal.Decimal) decimal.Decimal {
	if !n.IsInteger() || !p.IsInteger() {
		return decimal.Zero
	}
	nInt := n.IntPart()
	pInt := p.IntPart()
	return decimal.NewFromInt(int64(combin.NumPermutations(int(nInt), int(pInt))))
})

var StatPackage = ActionPackage{
	staticActions: []Action{
		&combOp, &permOp,
	},
}

// Stack package

var dupOp = NewStackOp("dup", 1, 2, func(elts ...Variable) []Variable {
	return []Variable{elts[0], elts[0]}
})

var dup2Op = NewStackOp("dup2", 2, 4, func(elts ...Variable) []Variable {
	return []Variable{elts[1], elts[0], elts[1], elts[0]}
})

var dropOp = NewStackOp("drop", 1, 0, func(elts ...Variable) []Variable {
	return []Variable{}
})

var drop2Op = NewStackOp("drop2", 2, 0, func(elts ...Variable) []Variable {
	return []Variable{}
})

var swapOp = NewStackOp("swap", 2, 2, func(elts ...Variable) []Variable {
	return []Variable{elts[1], elts[0]}
})

// rot, roll, pick

var dupNOp = NewRawStackOpWithCheck("dupn", 1, CheckFirstInt, func(system System, stack *Stack) error {
	n, err := stack.Pop()
	if err != nil {
		return err
	}
	stackElts, err := stack.PopN(int(n.asNumericVar().value.IntPart()))
	if err != nil {
		return err
	}
	stack.PushN(stackElts)
	return nil
})

var depthAct = NewRawStackOpWithCheck("depth", 0, CheckNoop, func(system System, stack *Stack) error {
	stack.Push(CreateNumericVariableFromInt(stack.Size()))
	return nil
})

var StackPackage = ActionPackage{
	staticActions: []Action{
		&dupOp,
		&dup2Op,
		&dropOp,
		&drop2Op,
		&swapOp,
		&dupNOp,
		&depthAct,
	},
}

var storeAct = NewActionDesc("sto",
	2,
	CheckGen([]Type{TYPE_ALG_EXPR, TYPE_GENERIC}),
	func(system System, stack *Stack) error {

		name, _ := stack.Pop()
		value, _ := stack.Pop()
		memory := system.Memory()
		rootFolder := memory.getRoot()
		_, err := memory.createVariable(name.asIdentifierVar().value, rootFolder, value)
		return err
	})

var loadAct = NewActionDesc("load",
	1,
	CheckGen([]Type{TYPE_ALG_EXPR}),
	func(system System, stack *Stack) error {

		variable, err := stack.Pop()
		if err != nil {
			return err
		}
		idAsString := variable.asIdentifierVar().value
		memory := system.Memory()
		rootFolder := memory.getRoot()
		for _, varName := range rootFolder.variables {
			if varName.name == idAsString {
				stack.Push(varName.value)
				break
			}
		}
		return nil
	})

var crdirAct = NewActionDesc("crdir", 1, CheckNoop, func(system System, stack *Stack) error {
	variable, err := stack.Pop()
	if err != nil {
		return err
	}
	folderNameAsStr := variable.asIdentifierVar().value
	memory := system.Memory()
	currentFolder := memory.getCurrentFolder()
	_, err = memory.createFolder(folderNameAsStr, currentFolder)
	return err
})

/*
var purgeAct = NewActionDesc("purge", 1, CheckNoop, func(system System, stack *Stack) error {
	return nil
})
*/

// UPDIR to go upper dir
// eval dirname => goto inside dirname
// { HOME TOTO } goes to TOTO
// RCL recall ?
// PURGE / PGDIR (only empty / delete recursively)

var MemoryPackage = ActionPackage{
	staticActions: []Action{
		&storeAct,
		&loadAct,
		&crdirAct,
	},
}

var MiscPackage = ActionPackage{
	staticActions: []Action{},
}
