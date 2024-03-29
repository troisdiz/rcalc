package rcalc

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"runtime"
	"strings"
	"testing"
	"troisdizaines.com/rcalc/rcalc/parser"
)

type LoggingParserListener struct {
	subListener parser.RcalcListener
	depth       int
}

func (l *LoggingParserListener) EnterInstructionSequence(c *parser.InstructionSequenceContext) {
	l.logMethodCalled()
	l.subListener.EnterInstructionSequence(c)
}

func (l *LoggingParserListener) ExitInstructionSequence(c *parser.InstructionSequenceContext) {
	l.logMethodCalled()
	l.subListener.ExitInstructionSequence(c)
}

var _ parser.RcalcListener = (*LoggingParserListener)(nil)

func (l *LoggingParserListener) spacesForDepth() string {
	result := ""
	for i := 0; i < l.depth*4; i++ {
		result += " "
	}
	return result
}

func (l *LoggingParserListener) getMethodCalled() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	// Extract simple function name
	lastDotPos := strings.LastIndex(frame.Function, ".")
	functionName := frame.Function[lastDotPos+1:]
	return functionName
}

func (l *LoggingParserListener) logMethodCalled() {
	functionName := l.getMethodCalled()
	l.logMethodCalledImplf(functionName, "")
}

func (l *LoggingParserListener) logMethodCalledf(format string, fArgs ...any) {
	functionName := l.getMethodCalled()
	l.logMethodCalledImplf(functionName, format, fArgs...)
}

func (l *LoggingParserListener) logMethodCalledImplf(functionName string, format string, fArgs ...any) {
	if functionName != "EnterEveryRule" && functionName != "ExitEveryRule" {
		// EnterEveryRule
		if strings.HasPrefix(functionName, "Enter") {
			GetLogger().Debugf("%s%s%s", l.spacesForDepth(), functionName, fmt.Sprintf(format, fArgs...))
			l.depth++
		} else if strings.HasPrefix(functionName, "Exit") {
			l.depth--
			GetLogger().Debugf("%s%s%s", l.spacesForDepth(), functionName, fmt.Sprintf(format, fArgs...))
		} else {
			GetLogger().Debugf("%s%s%s", l.spacesForDepth(), functionName, fmt.Sprintf(format, fArgs...))
		}
	}
}

func (l *LoggingParserListener) VisitTerminal(node antlr.TerminalNode) {
	l.logMethodCalledf(" => #%s# / #%d#", node.GetSymbol().GetText(), node.GetSymbol().GetTokenType())
	l.subListener.VisitTerminal(node)
}

func (l *LoggingParserListener) VisitErrorNode(node antlr.ErrorNode) {
	l.logMethodCalled()
	l.subListener.VisitErrorNode(node)
}

func (l *LoggingParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	l.logMethodCalled()
	l.subListener.EnterEveryRule(ctx)
}

func (l *LoggingParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	l.logMethodCalled()
	l.subListener.ExitEveryRule(ctx)
}

func (l *LoggingParserListener) EnterStart(c *parser.StartContext) {
	l.logMethodCalled()
	l.subListener.EnterStart(c)
}

func (l *LoggingParserListener) EnterProgram_declaration(c *parser.Program_declarationContext) {
	l.logMethodCalled()
	l.subListener.EnterProgram_declaration(c)
}

func (l *LoggingParserListener) ExitProgram_declaration(c *parser.Program_declarationContext) {
	l.logMethodCalled()
	l.subListener.ExitProgram_declaration(c)
}

func (l *LoggingParserListener) EnterInstrActionOrVarCall(c *parser.InstrActionOrVarCallContext) {
	l.logMethodCalled()
	l.subListener.EnterInstrActionOrVarCall(c)
}

func (l *LoggingParserListener) EnterInstrOp(c *parser.InstrOpContext) {
	l.logMethodCalled()
	l.subListener.EnterInstrOp(c)
}

func (l *LoggingParserListener) EnterInstrVariable(c *parser.InstrVariableContext) {
	l.logMethodCalled()
	l.subListener.EnterInstrVariable(c)
}

func (l *LoggingParserListener) EnterInstrIfThenElse(c *parser.InstrIfThenElseContext) {
	l.logMethodCalled()
	l.subListener.EnterInstrIfThenElse(c)
}

func (l *LoggingParserListener) EnterInstrStartNextLoop(c *parser.InstrStartNextLoopContext) {
	l.logMethodCalled()
	l.subListener.EnterInstrStartNextLoop(c)
}

func (l *LoggingParserListener) EnterInstrForNextLoop(c *parser.InstrForNextLoopContext) {
	l.logMethodCalled()
	l.subListener.EnterInstrForNextLoop(c)
}

func (l *LoggingParserListener) EnterInstrLocalVarCreation(c *parser.InstrLocalVarCreationContext) {
	l.logMethodCalled()
	l.subListener.EnterInstrLocalVarCreation(c)
}

func (l *LoggingParserListener) EnterOp(c *parser.OpContext) {
	l.logMethodCalled()
	l.subListener.EnterOp(c)
}

func (l *LoggingParserListener) EnterIf_then_else(c *parser.If_then_elseContext) {
	l.logMethodCalled()
	l.subListener.EnterIf_then_else(c)
}

func (l *LoggingParserListener) EnterStart_next_loop(c *parser.Start_next_loopContext) {
	l.logMethodCalled()
	l.subListener.EnterStart_next_loop(c)
}

func (l *LoggingParserListener) EnterFor_next_loop(c *parser.For_next_loopContext) {
	l.logMethodCalled()
	l.subListener.EnterFor_next_loop(c)
}

func (l *LoggingParserListener) EnterLocalVarCreation(c *parser.LocalVarCreationContext) {
	l.logMethodCalledf(" => %s", c.GetText())
	l.subListener.EnterLocalVarCreation(c)
}

func (l *LoggingParserListener) EnterDeclarationVariable(c *parser.DeclarationVariableContext) {
	l.logMethodCalled()
	l.subListener.EnterDeclarationVariable(c)
}

func (l *LoggingParserListener) EnterStatementLocalVarProgram(c *parser.StatementLocalVarProgramContext) {
	l.logMethodCalled()
	l.subListener.EnterStatementLocalVarProgram(c)
}

func (l *LoggingParserListener) EnterStatementLocalVarAlgebraicExpression(c *parser.StatementLocalVarAlgebraicExpressionContext) {
	l.logMethodCalled()
	l.subListener.EnterStatementLocalVarAlgebraicExpression(c)
}

func (l *LoggingParserListener) EnterVariableNumber(c *parser.VariableNumberContext) {
	l.logMethodCalled()
	l.subListener.EnterVariableNumber(c)
}

func (l *LoggingParserListener) EnterVariableAlgebraicExpression(c *parser.VariableAlgebraicExpressionContext) {
	l.logMethodCalled()
	l.subListener.EnterVariableAlgebraicExpression(c)
}

func (l *LoggingParserListener) EnterVariableProgramDeclaration(c *parser.VariableProgramDeclarationContext) {
	l.logMethodCalled()
	l.subListener.EnterVariableProgramDeclaration(c)
}

func (l *LoggingParserListener) EnterVariableList(c *parser.VariableListContext) {
	l.logMethodCalled()
	l.subListener.EnterVariableList(c)
}
func (l *LoggingParserListener) EnterList(c *parser.ListContext) {
	l.logMethodCalled()
	l.subListener.EnterList(c)
}

func (l *LoggingParserListener) ExitList(c *parser.ListContext) {
	l.logMethodCalled()
	l.subListener.ExitList(c)
}

func (l *LoggingParserListener) EnterListItem(c *parser.ListItemContext) {
	l.logMethodCalled()
}

func (l *LoggingParserListener) ExitListItem(c *parser.ListItemContext) {
	l.logMethodCalled()
}

func (l *LoggingParserListener) EnterVariableVector(c *parser.VariableVectorContext) {
	l.logMethodCalled()
	l.subListener.EnterVariableVector(c)
}

func (l *LoggingParserListener) EnterQuoted_algebraic_expression(c *parser.Quoted_algebraic_expressionContext) {
	l.logMethodCalled()
	l.subListener.EnterQuoted_algebraic_expression(c)
}

func (l *LoggingParserListener) EnterAlgExprAddSub(c *parser.AlgExprAddSubContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprAddSub(c)
}

func (l *LoggingParserListener) EnterAlgExprMulDiv(c *parser.AlgExprMulDivContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprMulDiv(c)
}

func (l *LoggingParserListener) EnterAlgExprPow(c *parser.AlgExprPowContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprPow(c)
}

func (l *LoggingParserListener) EnterAlgExprAddSignedAtom(c *parser.AlgExprAddSignedAtomContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprAddSignedAtom(c)
}

func (l *LoggingParserListener) EnterAlgExprSubSignedAtom(c *parser.AlgExprSubSignedAtomContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprSubSignedAtom(c)
}

func (l *LoggingParserListener) EnterAlgExprFuncAtom(c *parser.AlgExprFuncAtomContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprFuncAtom(c)
}

func (l *LoggingParserListener) EnterAlgExprAtom(c *parser.AlgExprAtomContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprAtom(c)
}

func (l *LoggingParserListener) EnterAlgExprNumber(c *parser.AlgExprNumberContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprNumber(c)
}

func (l *LoggingParserListener) EnterAlgExprVariable(c *parser.AlgExprVariableContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprVariable(c)
}

func (l *LoggingParserListener) EnterAlgExprParen(c *parser.AlgExprParenContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprParen(c)
}

func (l *LoggingParserListener) EnterAlg_variable(c *parser.Alg_variableContext) {
	l.logMethodCalled()
	l.subListener.EnterAlg_variable(c)
}

func (l *LoggingParserListener) EnterAlgExprFuncCall(c *parser.AlgExprFuncCallContext) {
	l.logMethodCalled()
	l.subListener.EnterAlgExprFuncCall(c)
}

func (l *LoggingParserListener) EnterVector(c *parser.VectorContext) {
	l.logMethodCalled()
	l.subListener.EnterVector(c)
}

func (l *LoggingParserListener) EnterAction_or_var_call(c *parser.Action_or_var_callContext) {
	l.logMethodCalled()
	l.subListener.EnterAction_or_var_call(c)
}

func (l *LoggingParserListener) ExitStart(c *parser.StartContext) {
	l.logMethodCalled()
	l.subListener.ExitStart(c)
}

func (l *LoggingParserListener) ExitInstrActionOrVarCall(c *parser.InstrActionOrVarCallContext) {
	l.logMethodCalled()
	l.subListener.ExitInstrActionOrVarCall(c)
}

func (l *LoggingParserListener) ExitInstrOp(c *parser.InstrOpContext) {
	l.logMethodCalled()
	l.subListener.ExitInstrOp(c)
}

func (l *LoggingParserListener) ExitInstrVariable(c *parser.InstrVariableContext) {
	l.logMethodCalled()
	l.subListener.ExitInstrVariable(c)
}

func (l *LoggingParserListener) ExitInstrIfThenElse(c *parser.InstrIfThenElseContext) {
	l.logMethodCalled()
	l.subListener.ExitInstrIfThenElse(c)
}

func (l *LoggingParserListener) ExitInstrStartNextLoop(c *parser.InstrStartNextLoopContext) {
	l.logMethodCalled()
	l.subListener.ExitInstrStartNextLoop(c)
}

func (l *LoggingParserListener) ExitInstrForNextLoop(c *parser.InstrForNextLoopContext) {
	l.logMethodCalled()
	l.subListener.ExitInstrForNextLoop(c)
}

func (l *LoggingParserListener) ExitInstrLocalVarCreation(c *parser.InstrLocalVarCreationContext) {
	l.logMethodCalled()
	l.subListener.ExitInstrLocalVarCreation(c)
}

func (l *LoggingParserListener) ExitOp(c *parser.OpContext) {
	l.logMethodCalled()
	l.subListener.ExitOp(c)
}

func (l *LoggingParserListener) ExitIf_then_else(c *parser.If_then_elseContext) {
	l.logMethodCalled()
	l.subListener.ExitIf_then_else(c)
}

func (l *LoggingParserListener) ExitStart_next_loop(c *parser.Start_next_loopContext) {
	l.logMethodCalled()
	l.subListener.ExitStart_next_loop(c)
}

func (l *LoggingParserListener) ExitFor_next_loop(c *parser.For_next_loopContext) {
	l.logMethodCalled()
	l.subListener.ExitFor_next_loop(c)
}

func (l *LoggingParserListener) ExitLocalVarCreation(c *parser.LocalVarCreationContext) {
	l.logMethodCalled()
	l.subListener.ExitLocalVarCreation(c)
}

func (l *LoggingParserListener) ExitDeclarationVariable(c *parser.DeclarationVariableContext) {
	l.logMethodCalled()
	l.subListener.ExitDeclarationVariable(c)
}

func (l *LoggingParserListener) ExitStatementLocalVarProgram(c *parser.StatementLocalVarProgramContext) {
	l.logMethodCalled()
	l.subListener.ExitStatementLocalVarProgram(c)
}

func (l *LoggingParserListener) ExitStatementLocalVarAlgebraicExpression(c *parser.StatementLocalVarAlgebraicExpressionContext) {
	l.logMethodCalled()
	l.subListener.ExitStatementLocalVarAlgebraicExpression(c)
}

func (l *LoggingParserListener) ExitVariableNumber(c *parser.VariableNumberContext) {
	l.logMethodCalled()
	l.subListener.ExitVariableNumber(c)
}

func (l *LoggingParserListener) ExitVariableAlgebraicExpression(c *parser.VariableAlgebraicExpressionContext) {
	l.logMethodCalled()
	l.subListener.ExitVariableAlgebraicExpression(c)
}

func (l *LoggingParserListener) ExitVariableProgramDeclaration(c *parser.VariableProgramDeclarationContext) {
	l.logMethodCalled()
	l.subListener.ExitVariableProgramDeclaration(c)
}

func (l *LoggingParserListener) ExitVariableList(c *parser.VariableListContext) {
	l.logMethodCalled()
	l.subListener.ExitVariableList(c)
}

func (l *LoggingParserListener) ExitVariableVector(c *parser.VariableVectorContext) {
	l.logMethodCalled()
	l.subListener.ExitVariableVector(c)
}

func (l *LoggingParserListener) ExitQuoted_algebraic_expression(c *parser.Quoted_algebraic_expressionContext) {
	l.logMethodCalled()
	l.subListener.ExitQuoted_algebraic_expression(c)
}

func (l *LoggingParserListener) ExitAlgExprAddSub(c *parser.AlgExprAddSubContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprAddSub(c)
}

func (l *LoggingParserListener) ExitAlgExprMulDiv(c *parser.AlgExprMulDivContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprMulDiv(c)
}

func (l *LoggingParserListener) ExitAlgExprPow(c *parser.AlgExprPowContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprPow(c)
}

func (l *LoggingParserListener) ExitAlgExprAddSignedAtom(c *parser.AlgExprAddSignedAtomContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprAddSignedAtom(c)
}

func (l *LoggingParserListener) ExitAlgExprSubSignedAtom(c *parser.AlgExprSubSignedAtomContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprSubSignedAtom(c)
}

func (l *LoggingParserListener) ExitAlgExprFuncAtom(c *parser.AlgExprFuncAtomContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprFuncAtom(c)
}

func (l *LoggingParserListener) ExitAlgExprAtom(c *parser.AlgExprAtomContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprAtom(c)
}

func (l *LoggingParserListener) ExitAlgExprNumber(c *parser.AlgExprNumberContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprNumber(c)
}

func (l *LoggingParserListener) ExitAlgExprVariable(c *parser.AlgExprVariableContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprVariable(c)
}

func (l *LoggingParserListener) ExitAlgExprParen(c *parser.AlgExprParenContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprParen(c)
}

func (l *LoggingParserListener) ExitAlg_variable(c *parser.Alg_variableContext) {
	l.logMethodCalled()
	l.subListener.ExitAlg_variable(c)
}

func (l *LoggingParserListener) ExitAlgExprFuncCall(c *parser.AlgExprFuncCallContext) {
	l.logMethodCalled()
	l.subListener.ExitAlgExprFuncCall(c)
}

func (l *LoggingParserListener) ExitVector(c *parser.VectorContext) {
	l.logMethodCalled()
	l.subListener.ExitVector(c)
}

func (l *LoggingParserListener) ExitAction_or_var_call(c *parser.Action_or_var_callContext) {
	l.logMethodCalled()
	l.subListener.ExitAction_or_var_call(c)
}

func (l *LoggingParserListener) EnterNumber(c *parser.NumberContext) {
	l.logMethodCalled()
	l.subListener.EnterNumber(c)
}

func (l *LoggingParserListener) ExitNumber(c *parser.NumberContext) {
	l.logMethodCalled()
	l.subListener.ExitNumber(c)
}

type ParsingTestSuite struct {
	suite.Suite

	registry *ActionRegistry
}

func TestParsingSuite(t *testing.T) {
	InitDevLogger("-")

	testingSuite := new(ParsingTestSuite)
	testingSuite.registry = initRegistry()
	suite.Run(t, testingSuite)
}

func (suite *ParsingTestSuite) parseWithDebugLogging(txt string) ([]Action, error) {
	return parseToActionsImpl(txt, "Test", suite.registry, func(listener parser.RcalcListener) parser.RcalcListener {
		return &LoggingParserListener{
			subListener: listener,
		}
	})
}

func (suite *ParsingTestSuite) TestDecimalFormats() {
	var strings = []string{"37", "4.5", "-0.4", "+.58", "-1e-12", "-.2e13"}
	for _, str := range strings {
		number, err := decimal.NewFromString(str)
		if assert.NoError(suite.T(), err, "could not parse: %s", str) {
			fmt.Printf("%s -> %v\n", str, number)
		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParse2Numbers() {
	var numbersToParse = []string{
		"37",
		"4.5",
		"-0.4",
		".58",
	}

	for _, expr := range numbersToParse {
		suite.Run(expr, func() {
			elt, err := suite.parseWithDebugLogging(expr)
			if assert.NoError(suite.T(), err, "Parse error : %s", err) {
				assert.IsType(suite.T(), elt[0], &VariablePutOnStackActionDesc{}, "type %t is not VariablePutOnStackActionDesc", elt[0])
			}
		})
	}
}

func (suite *ParsingTestSuite) TestAntlrIdentifierParser() {
	var txt string = "'ab' 'cd' 'de'"
	actions, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error: %s", err) {
		assert.Len(suite.T(), actions, 3)
	}
}

func (suite *ParsingTestSuite) TestAntlrAlgebraicExprParser() {
	var txt string = "'1+2'"
	actions, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error: %s", err) {
		assert.Len(suite.T(), actions, 1)

		assert.IsType(suite.T(), &VariablePutOnStackActionDesc{}, actions[0])

		actionDesc := actions[0].(*VariablePutOnStackActionDesc)
		assert.NotNil(suite.T(), actionDesc.value)
		assert.IsType(suite.T(), &AlgebraicExpressionVariable{}, actionDesc.value)

		algExprVar := actionDesc.value.(*AlgebraicExpressionVariable)
		assert.NotNil(suite.T(), algExprVar.rootNode)
		numericValue, _ := algExprVar.rootNode.Evaluate(nil)
		expected := decimal.NewFromInt(3)
		assert.Equal(suite.T(), expected, numericValue.value, "Expected %v / Value %v", expected, numericValue.value)
	}
}

func (suite *ParsingTestSuite) TestAntlrParseActionInRegistry() {
	var txt string = "quit sto"

	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 2) {
			assert.Equal(suite.T(), elt[0], &EXIT_ACTION)
		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParseStartNextLoop() {
	var txt string = "1 3 start 1 next"

	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 3) {
			assert.IsType(suite.T(), &StartNextLoopActionDesc{}, elt[2])
		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParseForNextLoop() {
	var txt string = "1 3 for i 1 next"

	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 3) {
			assert.IsType(suite.T(), &ForNextLoopActionDesc{}, elt[2])
			forNextLoopActionDesc := elt[2].(*ForNextLoopActionDesc)
			loopActions := forNextLoopActionDesc.actions
			if assert.Len(suite.T(), loopActions, 1) {
				assert.IsType(suite.T(), &VariablePutOnStackActionDesc{}, loopActions[0])
			}
		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParseForNextLoopError() {
	var txt string = "1 3 for i 1"

	_, err := suite.parseWithDebugLogging(txt)

	assert.Errorf(suite.T(), err, "")
}

func (suite *ParsingTestSuite) TestAntlrParseIfThenElse() {

	var txt string = " if 1 1 == then 2 else 3 end"

	GetLogger().Debugf("Parsing %s", txt)
	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 1) {
			assert.IsType(suite.T(), &IfThenElseActionDesc{}, elt[0])
			ifThenElseActionDesc := elt[0].(*IfThenElseActionDesc)
			ifActions := ifThenElseActionDesc.ifActions
			thenActions := ifThenElseActionDesc.thenActions
			elseActions := ifThenElseActionDesc.elseActions

			if assert.Len(suite.T(), ifActions, 3) {
				assert.IsType(suite.T(), &eqNumOp, ifActions[2])
			}
			if assert.Len(suite.T(), thenActions, 1) {
				assert.IsType(suite.T(), &VariablePutOnStackActionDesc{}, thenActions[0])
			}
			if assert.Len(suite.T(), elseActions, 1) {
				assert.IsType(suite.T(), &VariablePutOnStackActionDesc{}, elseActions[0])
			}

		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParseProgram() {

	var txt string = " << 1 3 for i 1 next >>"
	//var txt string = " << 1 >>"

	GetLogger().Debugf("Parsing %s", txt)
	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 1) {
			assert.IsType(suite.T(), &VariablePutOnStackActionDesc{}, elt[0])
			variablePutOnStackActionDesc := elt[0].(*VariablePutOnStackActionDesc)
			genericVariable := variablePutOnStackActionDesc.value
			if assert.IsType(suite.T(), &ProgramVariable{}, genericVariable) {
				programVariable := genericVariable.(*ProgramVariable)

				if assert.Len(suite.T(), programVariable.actions, 3) {
					assert.IsType(suite.T(), &ForNextLoopActionDesc{}, programVariable.actions[2])
				}
			}
		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParseLocalVariableDeclarationForProgram() {
	var txt string = " ->  a b << a >>"

	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 1) {
			assert.IsType(suite.T(), &VariableDeclarationActionDesc{}, elt[0])
			variableDeclarationActionDesc := elt[0].(*VariableDeclarationActionDesc)
			varNames := variableDeclarationActionDesc.varNames
			if assert.Len(suite.T(), varNames, 2) {
				assert.Equal(suite.T(), "a", varNames[0])
				assert.Equal(suite.T(), "b", varNames[1])
			}
			variable := variableDeclarationActionDesc.variableToEvaluate
			if assert.IsType(suite.T(), &ProgramVariable{}, variable) {
				programVariable := variable.asProgramVar()
				if assert.NotNil(suite.T(), programVariable) {
					if assert.Len(suite.T(), programVariable.actions, 1) {
						assert.IsType(suite.T(), &VariableEvaluationActionDesc{}, programVariable.actions[0])
					}
				}
			}
		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParseLocalVariableDeclarationForAlgebraicExpression() {
	var txt string = " -> a b 'a+b' "

	GetLogger().Debugf("Parsing %s", txt)
	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 1) {
			assert.IsType(suite.T(), &VariableDeclarationActionDesc{}, elt[0])
			variableDeclarationActionDesc := elt[0].(*VariableDeclarationActionDesc)
			varNames := variableDeclarationActionDesc.varNames
			if assert.Len(suite.T(), varNames, 2) {
				assert.Equal(suite.T(), "a", varNames[0])
				assert.Equal(suite.T(), "b", varNames[1])
			}
			variable := variableDeclarationActionDesc.variableToEvaluate
			if assert.IsType(suite.T(), &AlgebraicExpressionVariable{}, variable) {
				algExprVariable := variable.(*AlgebraicExpressionVariable)
				assert.NotNil(suite.T(), algExprVariable)
			}
		}
	}
}

func (suite *ParsingTestSuite) TestAntlrParseList() {
	var txt string = "{ 2 { 3 } }"

	elt, err := suite.parseWithDebugLogging(txt)
	if assert.NoError(suite.T(), err, "Parse error : %s", err) {
		fmt.Println(elt)
		if assert.Len(suite.T(), elt, 1) {
			assert.IsType(suite.T(), &VariablePutOnStackActionDesc{}, elt[0])
			variablePutListOnStack := elt[0].(*VariablePutOnStackActionDesc)
			listVar := variablePutListOnStack.value
			if assert.NotNil(suite.T(), listVar) {
				if assert.IsType(suite.T(), &ListVariable{}, listVar) {
					listVariable := listVar.(*ListVariable)
					assert.Len(suite.T(), listVariable.items, 2)
					assert.IsType(suite.T(), &ListVariable{}, listVariable.items[1])
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

func (t *TestErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (t *TestErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (t *TestErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

type TestParserListener struct {
	*parser.BaseRcalcListener
}

var _ parser.RcalcListener = (*TestParserListener)(nil)

func (suite *ParsingTestSuite) TestAlgebraicExpressionParsing() {

	expressions := []struct {
		literal string
		value   decimal.Decimal
	}{
		{
			literal: "'1 +2'",
			value:   decimal.NewFromInt(3),
		},
		{
			literal: "'1 + 2'", value: decimal.NewFromInt(3),
		},
		{
			literal: "'1 +2 - 5'", value: decimal.NewFromInt(-2),
		},
		{
			literal: "'1 +2 -5'", value: decimal.NewFromInt(-2),
		},
		{
			literal: "'1+ 2'",
			value:   decimal.NewFromInt(3),
		},
		{
			literal: "'1 * -2'",
			value:   decimal.NewFromInt(-2),
		},
		{
			literal: "'1 * +2'",
			value:   decimal.NewFromInt(2),
		},
		{
			literal: "'1*(2+ 3)'",
			value:   decimal.NewFromInt(5),
		},
		{
			literal: "'1*cos(2+3- 5)'",
			value:   decimal.NewFromInt(1),
		},
		{
			literal: "'-sin((2+3)*0)'",
			value:   decimal.NewFromInt(0),
		},
		{
			literal: "'1 + 2 + 3'",
			value:   decimal.NewFromInt(6),
		},
		{
			literal: "'1 + 2 - 3'",
			value:   decimal.Zero,
		},
		{
			literal: "'a'",
			value:   decimal.NewFromInt(7),
		},
		{
			literal: "'a+2'",
			value:   decimal.NewFromInt(9),
		},
		{
			literal: "'2^2'",
			value:   decimal.NewFromInt(4),
		},
		{
			literal: "'2^2^2'",
			value:   decimal.NewFromInt(16),
		},
		{
			literal: "'2^(2+3)'",
			value:   decimal.NewFromInt(32),
		},
	}

	//var nodeByExpression map[string]AlgebraicExpressionNode = make(map[string]AlgebraicExpressionNode)
	stack := CreateStack()
	system := CreateSystemInstance()
	_, err := system.memory.createVariable(
		"a",
		system.memory.getRoot(),
		CreateNumericVariable(decimal.NewFromInt(7)))
	var runtimeContext *RuntimeContext
	if assert.NoError(suite.T(), err, "Cannot create variable") {

		runtimeContext = CreateRuntimeContext(system, stack)

	}
	for idx, expr := range expressions {
		suite.Run(fmt.Sprintf("%02d-%s", idx+1, expr.literal), func() {

			var nodeExpr AlgebraicExpressionNode
			suite.T().Run("a - Parse", func(t *testing.T) {
				is := antlr.NewInputStream(expr.literal)
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
				antlr.ParseTreeWalkerDefault.Walk(listener, p.Start_())
				assert.False(t, el.hasErrors, "errors while parsing")
				expressionNodes, _ := listener.contextManager.actionCtxStack.GetCurrentRoot().CreateFinalItem()
				variablePutOnStackAction := expressionNodes[0].(*VariablePutOnStackActionDesc)
				algExprVariable := variablePutOnStackAction.value.(*AlgebraicExpressionVariable)
				if assert.NotNil(t, algExprVariable.rootNode, "Value of PutOnStackAction is nil for expr %s", expr.literal) {
					nodeExpr = algExprVariable.rootNode
				}
			})

			suite.T().Run("b - Compute", func(t *testing.T) {
				algExprNode := nodeExpr
				if algExprNode == nil {
					assert.Failf(t, "Parsing failed, cannot do compute test", "")
				}
				numericVariable, err := evalAlgExpression(runtimeContext, algExprNode)
				if assert.NoError(t, err) {
					assert.True(t,
						expr.value.Equal(numericVariable.value),
						"%s -> %v instead of %v\n", expr.literal, numericVariable.value, expr.value)
				}
			})
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
