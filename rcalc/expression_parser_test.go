package rcalc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseNumbers(t *testing.T) {
	var txt string = "3"
	var registry *ActionRegistry = initRegistry()

	lex := Lex("Test", txt)
	elt, err := ParseToActions(lex, registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
	}

	if assert.Len(t, elt, 1) {
		assert.IsType(t, elt[0], &VariablePutOnStackActionDesc{})
	}
}

func TestParse2Numbers(t *testing.T) {
	var txt string = "3 4.5"
	var registry *ActionRegistry = initRegistry()

	lex := Lex("Test", txt)
	elt, err := ParseToActions(lex, registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
	}

	if assert.Len(t, elt, 2) {
		assert.IsType(t, elt[0], &VariablePutOnStackActionDesc{})
		assert.IsType(t, elt[1], &VariablePutOnStackActionDesc{})
	}
}

func TestParseAddNumbers(t *testing.T) {
	var txt string = "3 4.5 + -"
	var registry *ActionRegistry = initRegistry()

	fmt.Printf("Text to parse: \"%s\"\n", txt)

	lex := Lex("Test", txt)
	elt, err := ParseToActions(lex, registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Printf("%v\n", elt)
	}

	if assert.Len(t, elt, 4) {
		assert.IsType(t, elt[0], &VariablePutOnStackActionDesc{})
		assert.IsType(t, elt[1], &VariablePutOnStackActionDesc{})
	}
}

func TestParseActionInRegistry(t *testing.T) {
	var txt string = "quit"
	var registry *ActionRegistry = initRegistry()

	lex := Lex("Test", txt)
	elt, err := ParseToActions(lex, registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
	}

	if assert.Len(t, elt, 1) {
		assert.Equal(t, elt[0], &EXIT_ACTION)
	}
}

func TestParseIdentifier(t *testing.T) {
	var txt string = "'ab' 'cd' 'de'"
	var registry *ActionRegistry = initRegistry()

	lex := Lex("Test", txt)
	elt, err := ParseToActions(lex, registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
	}

	if assert.Len(t, elt, 3) {
		assert.Equal(t, elt[0], &VariablePutOnStackActionDesc{
			value: CreateIdentifierVariable("ab"),
		})
	}
}
