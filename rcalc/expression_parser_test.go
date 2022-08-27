package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecimalFormats(t *testing.T) {
	var strings = []string{"37", "4.5", "-0.4", "+.58", "-1e-12", "-.2e13"}
	for _, str := range strings {
		number, err := decimal.NewFromString(str)
		if assert.NoError(t, err, "could not parse: %s", str) {
			fmt.Printf("%s -> %v\n", str, number)
		}
	}
}

func TestAntlrParse2Numbers(t *testing.T) {
	var txt string = "37 4.5 -0.4 +.58"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(t, elt, 4) {
			assert.IsType(t, elt[0], &VariablePutOnStackActionDesc{})
			assert.IsType(t, elt[1], &VariablePutOnStackActionDesc{})
			assert.IsType(t, elt[2], &VariablePutOnStackActionDesc{})
			assert.IsType(t, elt[3], &VariablePutOnStackActionDesc{})
		}
	}
}

func TestAntlrIdentifierParser(t *testing.T) {
	var txt string = "'ab' 'cd' 'de'"
	var registry *ActionRegistry = initRegistry()
	actions, err := ParseToActions(txt, "", registry)
	if assert.NoError(t, err, "Parse error: %s", err) {
		assert.Len(t, actions, 3)
	}
}

func TestAntlrParseActionInRegistry(t *testing.T) {
	var txt string = "quit sto"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(t, elt, 2) {
			assert.Equal(t, elt[0], &EXIT_ACTION)
		}
	}
}

func TestAntlrParseStartNextLoop(t *testing.T) {
	var txt string = "1 3 start 1 next"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(t, elt, 3) {
			assert.IsType(t, &StartNextLoopActionDesc{}, elt[2])
		}
	}
}

func TestAntlrParseForNextLoop(t *testing.T) {
	var txt string = "1 3 for i 1 next"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(t, elt, 3) {
			assert.IsType(t, &ForNextLoopActionDesc{}, elt[2])
			forNextLoopActionDesc := elt[2].(*ForNextLoopActionDesc)
			loopActions := forNextLoopActionDesc.actions
			if assert.Len(t, loopActions, 1) {
				assert.IsType(t, &VariablePutOnStackActionDesc{}, loopActions[0])
			}
		}
	}
}

func TestAntlrParseForNextLoopError(t *testing.T) {
	var txt string = "1 3 for i 1"
	var registry *ActionRegistry = initRegistry()

	_, err := ParseToActions(txt, "Test", registry)

	assert.Errorf(t, err, "")
}

func TestAntlrParseProgram(t *testing.T) {
	var txt string = " << 1 3 for i 1 next >>"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(t, elt, 1) {
			assert.IsType(t, &VariablePutOnStackActionDesc{}, elt[0])
			variablePutOnStackActionDesc := elt[0].(*VariablePutOnStackActionDesc)
			genericVariable := variablePutOnStackActionDesc.value
			if assert.IsType(t, &ProgramVariable{}, genericVariable) {
				programVariable := genericVariable.(*ProgramVariable)

				if assert.Len(t, programVariable.actions, 3) {
					assert.IsType(t, &ForNextLoopActionDesc{}, programVariable.actions[2])
				}
			}
		}
	}
}
