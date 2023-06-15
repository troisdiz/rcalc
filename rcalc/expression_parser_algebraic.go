package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"troisdizaines.com/rcalc/rcalc/parser"
)

type AlgebraicVariableContext struct {
	BaseParseContext[Variable] // to avoid reimplementing the interface
	parseContextManager        *ParseContextManager
	exprText                   string
}

var _ ParseContext[Variable] = (*AlgebraicVariableContext)(nil)

func (ac *AlgebraicVariableContext) CreateFinalItem() ([]Variable, error) {
	if len(ac.parseContextManager.lastAlgebraicValues) > 1 {
		return nil, fmt.Errorf("cannot create multiple variables with a single algebraic expression (should panic?)")
	}
	variable := &AlgebraicExpressionVariable{
		CommonVariable: CommonVariable{
			fType: TYPE_ALG_EXPR,
		},
		value:    ac.exprText,
		rootNode: ac.parseContextManager.lastAlgebraicValues[0],
	}
	// TODO must not be done here!!
	ac.parseContextManager.lastVariableValues = nil
	return []Variable{variable}, nil
}

func (ac *AlgebraicVariableContext) TokenVisited(token int) {

}

type RootAlgebraicExprContext struct {
	BaseParseContext[AlgebraicExpressionNode] // to avoid reimplementing the interface
}

func (rac *RootAlgebraicExprContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {
	//TODO test length
	return []AlgebraicExpressionNode{rac.items[0].item}, nil
}

var _ ParseContext[AlgebraicExpressionNode] = (*RootAlgebraicExprContext)(nil)

type AlgebraicExprContext struct {
	BaseParseContext[AlgebraicExpressionNode] // to avoid reimplementing the interface
	reg                                       *ActionRegistry
	tokens                                    []int
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicExprContext)(nil)

func (aec *AlgebraicExprContext) TokenVisited(token int) {
	aec.tokens = append(aec.tokens, token)
}

func (aec *AlgebraicExprContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {
	panic("AlgebraicExprContext.CreateFinalItem() must be override")
}

func (aec *AlgebraicExprContext) GetRootExprNode() AlgebraicExpressionNode {
	return aec.GetItems()[0].item
}

type AlgebraicAddSubContext struct {
	AlgebraicExprContext
}

func (asc *AlgebraicAddSubContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {

	items := asc.GetItems()
	if len(items) == 1 {
		return []AlgebraicExpressionNode{items[0].item}, nil
	}
	operators, err := tokenToPosition([]int{parser.RcalcLexerOP_ADD, parser.RcalcLexerOP_SUB}, asc.tokens)
	fmt.Printf("Length of operators is %d / tokens : %d\n", len(operators), len(asc.tokens))
	if err != nil {
		panic("Unknown token")
	}
	return []AlgebraicExpressionNode{
		&AlgExprAddSub{
			items:     toNonLocated(items),
			operators: operators,
		},
	}, nil
}

type AlgebraicMulDivContext struct {
	AlgebraicExprContext
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicMulDivContext)(nil)

func (amdc *AlgebraicMulDivContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {
	items := amdc.GetItems()
	if len(items) == 1 {
		return []AlgebraicExpressionNode{items[0].item}, nil
	}
	operators, err := tokenToPosition([]int{parser.RcalcLexerOP_MUL, parser.RcalcLexerOP_DIV}, amdc.tokens)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic("Unknown token")
	}

	return []AlgebraicExpressionNode{
		&AlgExprMulDiv{
			items:     toNonLocated(amdc.GetItems()),
			operators: operators,
		},
	}, nil
}

type AlgebraicPowerContext struct {
	AlgebraicExprContext
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicPowerContext)(nil)

func (amdc *AlgebraicPowerContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {
	items := amdc.GetItems()
	if len(items) == 1 {
		return []AlgebraicExpressionNode{items[0].item}, nil
	}
	_, err := tokenToPosition([]int{parser.RcalcLexerOP_POW}, amdc.tokens)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic("Unknown token")
	}

	return []AlgebraicExpressionNode{
		&AlgExprPow{
			items: toNonLocated(amdc.GetItems()),
		},
	}, nil
}

type AlgebraicSignedAtomContext struct {
	AlgebraicExprContext
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicSignedAtomContext)(nil)

func (asac *AlgebraicSignedAtomContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {
	operators, err := tokenToPosition([]int{parser.RcalcLexerOP_ADD, parser.RcalcLexerOP_SUB}, asac.tokens)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic("Unknown token")
	}

	return []AlgebraicExpressionNode{
		&AlgExprSignedElt{
			items:    asac.GetItems()[0].item,
			operator: operators[0],
		},
	}, nil
}

type AlgebraicFunctionContext struct {
	AlgebraicExprContext

	functionName string
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicFunctionContext)(nil)

func (afc *AlgebraicFunctionContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {

	if fn := afc.reg.GetAlgebraicFunction(afc.functionName); fn != nil {
		return []AlgebraicExpressionNode{
			&AlgExprFunctionElt{
				functionName: afc.functionName,
				fn:           fn,
				arguments:    toNonLocated(afc.GetItems()),
			},
		}, nil
	} else {
		// TODO error handling of such cases
		return nil, nil
	}
}

type AlgebraicAtomContext struct {
	AlgebraicExprContext
}

func (aac *AlgebraicAtomContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {

	operator := OPERATOR_ADD
	if len(aac.tokens) > 0 && aac.tokens[0] == parser.RcalcLexerOP_SUB {
		operator = OPERATOR_SUB
	}
	return []AlgebraicExpressionNode{
		&AlgExprSignedElt{
			items:    aac.GetItems()[0].item,
			operator: operator,
		},
	}, nil
}

type AlgebraicNumberContext struct {
	AlgebraicExprContext

	value decimal.Decimal
}

func (anc *AlgebraicNumberContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {

	return []AlgebraicExpressionNode{
		&AlgExprNumber{
			value: anc.value,
		},
	}, nil
}

type AlgebraicVariableNameContext struct {
	AlgebraicExprContext

	value string
}

func (avnc *AlgebraicVariableNameContext) CreateFinalItem() ([]AlgebraicExpressionNode, error) {
	return []AlgebraicExpressionNode{
		&AlgExprVariable{
			value: avnc.value,
		},
	}, nil
}
