package rcalc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseNumbers(t *testing.T) {
	var txt string = "3"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions(txt, "Test", registry)
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

	elt, err := ParseToActions(txt, "Test", registry)
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

	elt, err := ParseToActions(txt, "Test", registry)
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

	elt, err := ParseToActions(txt, "Test", registry)
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

	elt, err := ParseToActions(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
	}

	if assert.Len(t, elt, 3) {
		assert.Equal(t, elt[0], &VariablePutOnStackActionDesc{
			value: CreateIdentifierVariable("ab"),
		})
	}
}

func TestAntlrIdentifierParser(t *testing.T) {
	var txt string = "'ab' 'cd' 'de'"
	var registry *ActionRegistry = initRegistry()
	actions, err := ParseToActions2(txt, "", registry)
	if assert.NoError(t, err, "Parse error: %s", err) {
		assert.Len(t, actions, 3)
	}
}

func TestAntlrParseActionInRegistry(t *testing.T) {
	var txt string = "quit sto"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions2(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(t, elt, 2) {
			assert.Equal(t, elt[0], &EXIT_ACTION)
		}
	}
}
