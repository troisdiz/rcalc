package rcalc

import (
    "fmt"
    "strconv"
    "strings"
)

type ExprElementType int
const (
    ACTION_EXPR_TYPE ExprElementType = 0
    STACK_ELT_EXPR_TYPE ExprElementType= 1
)

type ExprElement struct {
    eltType ExprElementType
    elt interface{}
}

func (e *ExprElement) String() string {
    var eltStr string
    switch e.eltType {
    case ACTION_EXPR_TYPE:
        eltStr = fmt.Sprint(e.elt.(*ActionDesc))
    case STACK_ELT_EXPR_TYPE:
        eltStr = fmt.Sprint(e.elt.(StackElt))
    }
    return fmt.Sprintf("ExprElt(%d, %s)", e.eltType, eltStr)
}

func createActionExprElt(action *ActionDesc) *ExprElement {
    return &ExprElement{
        eltType: ACTION_EXPR_TYPE,
        elt: action,
    }
}

func (e *ExprElement) asAction() Action {
    return e.elt.(Action)
}

func createStackEltExprElt(stackElt StackElt) *ExprElement {
    return &ExprElement{
        eltType: STACK_ELT_EXPR_TYPE,
        elt: stackElt,
    }
}

func (e *ExprElement) asStackElt() StackElt {
    return e.elt.(StackElt)
}

func ParseExpression(registry *ActionRegistry, input string) ([]*ExprElement, error) {
    var result []*ExprElement

    for _, exprEltStr := range strings.Split(input, " ") {
        trimmedElt := strings.TrimSpace(exprEltStr)
        parsedElt, err := parseExpressionElt(registry, trimmedElt)
        if err != nil {
            return nil, err
        } else {
            result = append(result, parsedElt)
        }
    }
    return result, nil
}


func parseExpressionElt(registry *ActionRegistry, elt string) (*ExprElement,error) {
    if registry.ContainsOpCode(elt) {
        action := registry.GetAction(elt)
        return createActionExprElt(action), nil
    }

    if value, err := strconv.Atoi(elt); err == nil {
        return createStackEltExprElt(CreateIntStackElt(value)), nil
    }
    return nil, fmt.Errorf("could not parse \"%s\"", elt)
}
