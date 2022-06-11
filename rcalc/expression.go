package rcalc

import (
	"fmt"
)

type ExprElementType int

const (
	ACTION_EXPR_TYPE    ExprElementType = 0
	STACK_ELT_EXPR_TYPE ExprElementType = 1
)

type ExprElement struct {
	eltType ExprElementType
	elt     interface{}
}

func (e *ExprElement) String() string {
	var eltStr string
	switch e.eltType {
	case ACTION_EXPR_TYPE:
		eltStr = fmt.Sprint(e.elt.(*OperationDesc))
	case STACK_ELT_EXPR_TYPE:
		eltStr = fmt.Sprint(e.elt.(StackElt))
	}
	return fmt.Sprintf("ExprElt(%d, %s)", e.eltType, eltStr)
}
