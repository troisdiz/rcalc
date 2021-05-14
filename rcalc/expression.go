package rcalc

import (
    "strconv"
    "strings"
)
/*

 */
type ExprElementType int
const (
    ACTION_EXPR_TYPE ExprElementType = 0
    STACK_ELT_EXPR_TYPE ExprElementType= 1
    OP_EXPR_TYPE ExprElementType= 2
)

type ExprElement struct {
    eltType ExprElementType
    elt interface{}
}

func createActionExprElt(action Action) *ExprElement {
    return &ExprElement{
        eltType: ACTION_EXPR_TYPE,
        elt: action,
    }
}

func (e *ExprElement) asAction() Action {
    return e.elt.(Action)
}

func createStackEltExprElt(stackElt StackElt) ExprElement {
    return ExprElement{
        eltType: STACK_ELT_EXPR_TYPE,
        elt: stackElt,
    }
}

func (e *ExprElement) asStackElt() StackElt {
    return e.elt.(StackElt)
}

func createOperationExprElt(op Operation) *ExprElement {
    return &ExprElement{
        eltType: OP_EXPR_TYPE,
        elt: op,
    }
}

func (e *ExprElement) asOp() Operation {
    return e.elt.(Operation)
}

func ParseExpression(input string) []ExprElement {
    var result []ExprElement

    for _, exprEltStr := range strings.Split(input, " ") {
        if value, err := strconv.Atoi(exprEltStr); err == nil {
            result = append(result, createStackEltExprElt(CreateInStackElt(value)))
        }
    }

    return result
}
