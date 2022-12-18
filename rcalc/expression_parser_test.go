package rcalc

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"troisdizaines.com/rcalc/rcalc/parser"
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
	var numbersToParse = []string{
		"37",
		"4.5",
		"-0.4",
		".58",
	}
	var registry *ActionRegistry = initRegistry()

	for _, expr := range numbersToParse {
		t.Run(expr, func(t *testing.T) {
			elt, err := ParseToActions(expr, "Test", registry)
			assert.NoError(t, err, "Parse error : %s", err)
			assert.IsType(t, elt[0], &VariablePutOnStackActionDesc{})
		})
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

func TestAntlrAlgebraicExprParser(t *testing.T) {
	var txt string = "'1+2'"
	var registry *ActionRegistry = initRegistry()
	actions, err := ParseToActions(txt, "", registry)
	if assert.NoError(t, err, "Parse error: %s", err) {
		assert.Len(t, actions, 1)

		assert.IsType(t, &VariablePutOnStackActionDesc{}, actions[0])

		actionDesc := actions[0].(*VariablePutOnStackActionDesc)
		assert.NotNil(t, actionDesc.value)
		assert.IsType(t, &AlgebraicExpressionVariable{}, actionDesc.value)

		algExprVar := actionDesc.value.(*AlgebraicExpressionVariable)
		assert.NotNil(t, algExprVar.rootNode)
		numericValue := algExprVar.rootNode.evaluate(nil)
		expected := decimal.NewFromInt(3)
		assert.Equal(t, expected, numericValue.value, "Expected %v / Value %v", expected, numericValue.value)
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

func TestAntlrParseIfThenElse(t *testing.T) {
	var txt string = " if 1 1 == then 2 else 3 end"
	var registry *ActionRegistry = initRegistry()

	elt, err := ParseToActions(txt, "Test", registry)
	if assert.NoError(t, err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(t, elt, 1) {
			assert.IsType(t, &IfThenElseActionDesc{}, elt[0])
			ifThenElseActionDesc := elt[0].(*IfThenElseActionDesc)
			ifActions := ifThenElseActionDesc.ifActions
			thenActions := ifThenElseActionDesc.thenActions
			elseActions := ifThenElseActionDesc.elseActions

			if assert.Len(t, ifActions, 3) {
				assert.IsType(t, &eqNumOp, ifActions[2])
			}
			if assert.Len(t, thenActions, 1) {
				assert.IsType(t, &VariablePutOnStackActionDesc{}, thenActions[0])
			}
			if assert.Len(t, elseActions, 1) {
				assert.IsType(t, &VariablePutOnStackActionDesc{}, elseActions[0])
			}

		}
	}
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

type TestErrorListener struct {
	hasErrors bool
}

var _ antlr.ErrorListener = (*TestErrorListener)(nil)

func (t *TestErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	t.hasErrors = true
}

func (t *TestErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

func (t *TestErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

func (t *TestErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
}

type TestParserListener struct {
	*parser.BaseRcalcListener
}

var _ parser.RcalcListener = (*TestParserListener)(nil)

func TestAlgebraicExpressionParsing(t *testing.T) {

	expressions := []string{
		"'1 +2'",
		"'1 + 2'",
		"'1+ 2'",
		"'1 * -2'",
		"'1 * +2'",
		"'1*(2+ 3)'",
		"'1*cos(2+3)'",
		"'-sin(2+3)'",
		"'1 + 2 + 3'",
		"'1 + 2 - 3'",
	}

	for idx, expr := range expressions {
		t.Run(fmt.Sprintf("Parse %02d-%s", idx+1, expr), func(t *testing.T) {
			is := antlr.NewInputStream(expr)
			// Create the Lexer
			lexer := parser.NewRcalcLexer(is)
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

			// Create the Parser
			p := parser.NewRcalcParser(stream)

			// Error Listener
			el := &TestErrorListener{}

			// Finally parse the expression (by walking the tree)
			var listener = CreateRcalcParserListener(Registry)
			//p.RemoveErrorListeners()
			p.AddErrorListener(el)
			antlr.ParseTreeWalkerDefault.Walk(listener, p.Start())
			assert.False(t, el.hasErrors)
			listener.rootPc.GetActions()
		})
	}
}

func TestTokenToPosition(t *testing.T) {
	ref := []int{3, 5, 12}
	tokens := []int{5, 3, 12}
	positions, err := tokenToPosition(ref, tokens)
	if assert.NoError(t, err) {
		if assert.Len(t, positions, 3) {
			assert.Equal(t, []int{1, 0, 2}, positions)
		}

	}
}
