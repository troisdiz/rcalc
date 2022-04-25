package rcalc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStackEltExpr(t *testing.T) {
	var s string = "3"
	var registry *ActionRegistry = initRegistry()
	elt, err := parseExpressionElt(registry, s)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
	}
}

func TestParseActionExpr(t *testing.T) {
	var s string = "quit"
	var registry *ActionRegistry = initRegistry()
	elt, err := parseExpressionElt(registry, s)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
	}
}

func TestParseAddition(t *testing.T) {
	var s string = "2 3 +"
	var registry *ActionRegistry = initRegistry()
	elts, err := ParseExpression(registry, s)
	if assert.NoError(t, err, "Parse error : %s", err) {
		for _, elt := range elts {
			fmt.Printf("%s\n", elt)
		}
	}
}
