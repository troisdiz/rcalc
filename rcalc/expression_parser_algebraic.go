package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"troisdizaines.com/rcalc/rcalc/parser"
)

type AlgebraicVariableContext struct {
	BaseParseContext[Action] // to avoid reimplementing the interface

	exprText         string
	algebraicContext *AlgebraicExprContext
}

func (ac *AlgebraicVariableContext) CreateFinalAction() Action {

	var algRootNode AlgebraicExpressionNode
	// TODO
	if ac.algebraicContext != nil {
		algRootNode = ac.algebraicContext.GetRootExprNode()
	}

	//AlgebraicExpressionNode{}
	return &VariablePutOnStackActionDesc{value: &AlgebraicExpressionVariable{
		CommonVariable: CommonVariable{
			fType: TYPE_ALG_EXPR,
		},
		value:    ac.exprText,
		rootNode: algRootNode,
	}}
}

func (ac *AlgebraicVariableContext) TokenVisited(token int) {

}

type AlgebraicExprContext struct {
	BaseParseContext[AlgebraicExpressionNode] // to avoid reimplementing the interface

	tokens []int
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicExprContext)(nil)

func (aec *AlgebraicExprContext) TokenVisited(token int) {
	aec.tokens = append(aec.tokens, token)
}

func (aec *AlgebraicExprContext) CreateFinalAction() AlgebraicExpressionNode {
	panic("AlgebraicExprContext.CreateFinalAction() must be override")
}

func (aec *AlgebraicExprContext) GetRootExprNode() AlgebraicExpressionNode {
	return aec.GetItems()[0]
}

type AlgebraicAddSubContext struct {
	AlgebraicExprContext
}

func (asc *AlgebraicAddSubContext) CreateFinalAction() AlgebraicExpressionNode {
	operators, err := tokenToPosition([]int{parser.RcalcLexerOP_ADD, parser.RcalcLexerOP_SUB}, asc.tokens)
	fmt.Printf("Length of operators is %d / tokens : %d\n", len(operators), len(asc.tokens))
	if err != nil {
		panic("Unknown token")
	}
	return &AlgExprAddSub{
		items:     asc.GetItems(),
		operators: operators,
	}
}

type AlgebraicMulDivContext struct {
	AlgebraicExprContext
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicMulDivContext)(nil)

func (amdc *AlgebraicMulDivContext) CreateFinalAction() AlgebraicExpressionNode {
	operators, err := tokenToPosition([]int{parser.RcalcLexerOP_MUL, parser.RcalcLexerOP_DIV}, amdc.tokens)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic("Unknown token")
	}

	return &AlgExprMulDiv{
		items:     amdc.GetItems(),
		operators: operators,
	}
}

type AlgebraicSignedAtomContext struct {
	AlgebraicExprContext
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicSignedAtomContext)(nil)

func (asac *AlgebraicSignedAtomContext) CreateFinalAction() AlgebraicExpressionNode {
	operators, err := tokenToPosition([]int{parser.RcalcLexerOP_ADD, parser.RcalcLexerOP_SUB}, asac.tokens)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic("Unknown token")
	}

	return &AlgExprSignedElt{
		items:    asac.GetItems()[0],
		operator: operators[0],
	}
}

type AlgebraicFunctionContext struct {
	AlgebraicExprContext

	functionName string
}

var _ ParseContext[AlgebraicExpressionNode] = (*AlgebraicFunctionContext)(nil)

func (afc *AlgebraicFunctionContext) CreateFinalAction() AlgebraicExpressionNode {

	return &AlgExprFunctionElt{
		functionName: afc.functionName,
		arguments:    afc.GetItems(),
	}
}

type AlgebraicAtomContext struct {
	AlgebraicExprContext
}

func (aac *AlgebraicAtomContext) CreateFinalAction() AlgebraicExpressionNode {

	operator := OPERATOR_ADD
	if len(aac.tokens) > 0 && aac.tokens[0] == parser.RcalcLexerOP_SUB {
		operator = OPERATOR_SUB
	}
	return &AlgExprSignedElt{
		items:    aac.GetItems()[0],
		operator: operator,
	}
}

type AlgebraicNumberContext struct {
	AlgebraicExprContext

	value decimal.Decimal
}

func (anc *AlgebraicNumberContext) CreateFinalAction() AlgebraicExpressionNode {

	return &AlgExprNumber{
		value: anc.value,
	}
}

type AlgebraicVariableNameContext struct {
	AlgebraicExprContext

	value string
}

func (avnc *AlgebraicVariableNameContext) CreateFinalAction() AlgebraicExpressionNode {
	return &AlgExprVariable{
		value: avnc.value,
	}
}
