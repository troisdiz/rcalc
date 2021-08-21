package rcalc

import (
    "fmt"
    "testing"
)

func TestParseStackEltExpr(t *testing.T) {
    var s string = "3"
    var registry *ActionRegistry = initRegistry()
    elt, err := parseExpressionElt(registry, s)
    if err != nil {
        t.Errorf("Parse error : %s", err)
    } else {
        fmt.Println(elt)
    }
}

func TestParseActionExpr(t *testing.T) {
    var s string = "quit"
    var registry *ActionRegistry = initRegistry()
    elt, err := parseExpressionElt(registry, s)
    if err != nil {
        t.Errorf("Parse error : %s", err)
    } else {
        fmt.Println(elt)
    }
}

func TestParseAddition(t *testing.T) {
    var s string = "2 3 +"
    var registry *ActionRegistry = initRegistry()
    elts, err := ParseExpression(registry, s)
    if err != nil {
        t.Errorf("Parse error : %s", err)
    } else {
        for _, elt := range elts {
            fmt.Printf("%s\n", elt)
        }
    }
}
