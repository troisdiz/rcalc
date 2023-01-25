package rcalc

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/shopspring/decimal"
	"strings"
	parser "troisdizaines.com/rcalc/rcalc/parser"
)

func tokenToPosition(ref []int, tokens []int) ([]int, error) {
	result := make([]int, len(tokens))
	var found bool
	for idx, token := range tokens {
		found = false
		for pos, candidate := range ref {
			if candidate == token {
				found = true
				result[idx] = pos
				continue
			}
		}
		if !found {
			//myLexer.getVocabulary.getSymbolicName(myTerminalNode.getSymbol.getType)

			return nil, fmt.Errorf("token %d at position %d is not expected", token, idx)
		}
	}
	return result, nil
}

func parserNumber(txt string) (Variable, error) {
	number, err := decimal.NewFromString(txt)
	if err != nil {
		return nil, err
	}
	return CreateNumericVariable(number), nil
}

func parseIdentifier(txt string, node AlgebraicExpressionNode) (Variable, error) {
	l := len(txt)
	// TODO is identifier also ?
	return CreateAlgebraicExpressionVariable(txt[1:l-1], node), nil
}

func parseAction(txt string, registry *ActionRegistry) (Action, error) {
	lowerTxt := strings.ToLower(txt)
	if registry.ContainsOpCode(lowerTxt) {
		return registry.GetAction(txt), nil
	} else {
		return nil, fmt.Errorf("unknown action")
	}
}

type Location struct {
	start, stop antlr.Token
}

type ValidationError struct {
	location Location
	err      error
}

type ParseContext[T any] interface {
	GetParent() ParseContext[T]
	SetParent(ctx ParseContext[T])

	AddItem(item T)
	AddIdentifier(id string)
	ReportValidationError(location Location, err error)
	GetValidationErrors() []ValidationError

	BackFromChild(child ParseContext[T])
	CreateFinalAction() (T, error)

	TokenVisited(token int)
}

type BaseParseContext[T any] struct {
	parent           ParseContext[T]
	items            []T
	idDeclarations   []string
	validationErrors []ValidationError
}

var _ ParseContext[string] = (*BaseParseContext[string])(nil)

func (g *BaseParseContext[T]) GetParent() ParseContext[T] {
	return g.parent
}

func (g *BaseParseContext[T]) SetParent(ctx ParseContext[T]) {
	g.parent = ctx
}

func (g *BaseParseContext[T]) AddItem(item T) {
	g.items = append(g.items, item)
}

func (g *BaseParseContext[T]) AddIdentifier(id string) {
	g.idDeclarations = append(g.idDeclarations, id)
}

func (g *BaseParseContext[T]) ReportValidationError(location Location, err error) {
	g.validationErrors = append(g.validationErrors, ValidationError{location: location, err: err})
}

func (g *BaseParseContext[T]) GetValidationErrors() []ValidationError {
	return g.validationErrors
}

func (g *BaseParseContext[T]) BackFromChild(child ParseContext[T]) {
	childValidationErrors := child.GetValidationErrors()
	if len(childValidationErrors) > 0 {
		for _, ve := range childValidationErrors {
			g.validationErrors = append(g.validationErrors, ve)
		}
	} else {
		action, err := child.CreateFinalAction()
		if err != nil {
			g.ReportValidationError(Location{}, err)
		} else {
			g.AddItem(action)
		}
	}
}

func (g *BaseParseContext[T]) CreateFinalAction() (T, error) {
	panic("CreateFinalAction must be implemented by sub structures")
}

func (g *BaseParseContext[T]) GetItems() []T {
	return g.items
}

func (g *BaseParseContext[T]) TokenVisited(token int) {}

type IfThenElseContext struct {
	BaseParseContext[Action] // to avoid reimplementing the interface

	actions       [][]Action
	currentAction int
}

func (i *IfThenElseContext) AddItem(action Action) {
	i.actions[i.currentAction] = append(i.actions[i.currentAction], action)
}

func (i *IfThenElseContext) TokenVisited(token int) {
	if i.actions == nil {
		i.actions = make([][]Action, 3)
	}
	switch token {
	case parser.RcalcLexerKW_IF:
		i.currentAction = 0
	case parser.RcalcLexerKW_THEN:
		i.currentAction = 1
	case parser.RcalcLexerKW_ELSE:
		i.currentAction = 2
	}
}

func (i *IfThenElseContext) CreateFinalAction() (Action, error) {
	return &IfThenElseActionDesc{
		ifActions:   i.actions[0],
		thenActions: i.actions[1],
		elseActions: i.actions[2],
	}, nil
}

type StartEndLoopContext struct {
	BaseParseContext[Action]
}

func (pc *StartEndLoopContext) CreateFinalAction() (Action, error) {
	return &StartNextLoopActionDesc{actions: pc.BaseParseContext.items}, nil
}

type ForNextLoopContext struct {
	BaseParseContext[Action]
}

func (pc *ForNextLoopContext) CreateFinalAction() (Action, error) {
	return &ForNextLoopActionDesc{
		varName: pc.BaseParseContext.idDeclarations[0],
		actions: pc.BaseParseContext.items,
	}, nil
}

type ProgramContext struct {
	BaseParseContext[Action]
}

func (pc *ProgramContext) CreateFinalAction() (Action, error) {
	progVar := CreateProgramVariable(pc.items)
	return &VariablePutOnStackActionDesc{value: progVar}, nil
}

type InstrLocalVarCreationContext struct {
	BaseParseContext[Action]
}

func (pc *InstrLocalVarCreationContext) CreateFinalAction() (Action, error) {
	programPutOnStackVariable := pc.BaseParseContext.items[0].(*VariablePutOnStackActionDesc)
	//fmt.Printf("%v\n", programPutOnStackVariable)
	programVariable := programPutOnStackVariable.value.asProgramVar()
	return &VariableDeclarationActionDesc{
		varNames:        pc.BaseParseContext.idDeclarations,
		programVariable: programVariable,
	}, nil
}

type RcalcParserListener struct {
	*parser.BaseRcalcListener

	registry *ActionRegistry

	rootPc    *BaseParseContext[Action]
	currentPc ParseContext[Action]

	rootAlgebraicPc    *AlgebraicExprContext
	currentAlgebraicPc ParseContext[AlgebraicExpressionNode]
}

func CreateRcalcParserListener(registry *ActionRegistry) *RcalcParserListener {
	rootPc := &BaseParseContext[Action]{
		parent: nil,
		items:  nil,
	}
	return &RcalcParserListener{
		registry:  registry,
		rootPc:    rootPc,
		currentPc: rootPc,
	}
}

func (l *RcalcParserListener) AddAction(action Action) {
	l.currentPc.AddItem(action)
}

func (l *RcalcParserListener) AddVarName(varName string) {
	l.currentPc.AddIdentifier(varName)
}

func (l *RcalcParserListener) StartNewContext(ctx ParseContext[Action]) {
	ctx.SetParent(l.currentPc)
	l.currentPc = ctx
}

func (l *RcalcParserListener) BackToParentContext() {
	oldCurrent := l.currentPc
	l.currentPc = l.currentPc.GetParent()
	l.currentPc.BackFromChild(oldCurrent)
}

func (l *RcalcParserListener) StartNewAlgebraicContext(ctx ParseContext[AlgebraicExpressionNode]) {
	ctx.SetParent(l.currentAlgebraicPc)
	l.currentAlgebraicPc = ctx
}

func (l *RcalcParserListener) BackToParentAlgebraicContext() {
	oldCurrent := l.currentAlgebraicPc
	l.currentAlgebraicPc = l.currentAlgebraicPc.GetParent()
	l.currentAlgebraicPc.BackFromChild(oldCurrent)
}

func (l *RcalcParserListener) TokenVisited(token int) {
	l.currentPc.TokenVisited(token)
	if l.currentAlgebraicPc != nil {
		l.currentAlgebraicPc.TokenVisited(token)
	}
}

// ExitVariableNumber is called when production InstrNumber is exited.
func (l *RcalcParserListener) ExitVariableNumber(ctx *parser.VariableNumberContext) {
	//fmt.Printf("ExitInstrNumber: %s\n", ctx.GetText())
	number, err := parserNumber(ctx.GetText())
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(&VariablePutOnStackActionDesc{number})
	}

}

// EnterVariableAlgebraicExpression is called when production VariableAlgebraicExpression is entered.
func (l *RcalcParserListener) EnterVariableAlgebraicExpression(ctx *parser.VariableAlgebraicExpressionContext) {
	//fmt.Println("EnterVariableAlgebraicExpression")

	l.rootAlgebraicPc = &AlgebraicExprContext{
		reg: l.registry,
	}
	l.currentAlgebraicPc = l.rootAlgebraicPc

	text := ctx.GetText()
	lenText := len(text)
	exprText := text[1 : lenText-1]
	l.StartNewContext(&AlgebraicVariableContext{
		exprText:         exprText,
		algebraicContext: l.rootAlgebraicPc,
	})
}

// ExitVariableAlgebraicExpression is called when production VariableAlgebraicExpression is exited.
func (l *RcalcParserListener) ExitVariableAlgebraicExpression(ctx *parser.VariableAlgebraicExpressionContext) {
	//fmt.Println("ExitVariableAlgebraicExpression")

	rootAlgExpr := l.rootAlgebraicPc.GetItems()

	identifier, err := parseIdentifier(ctx.GetText(), rootAlgExpr[0])
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(&VariablePutOnStackActionDesc{value: identifier})
	}
	l.BackToParentContext()
}

// EnterAlgExprAddSub is called when entering the AlgExprAddSub production.
func (l *RcalcParserListener) EnterAlgExprAddSub(c *parser.AlgExprAddSubContext) {
	l.StartNewAlgebraicContext(&AlgebraicAddSubContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprAddSub is called when production AlgExprRoot is exited.
func (l *RcalcParserListener) ExitAlgExprAddSub(ctx *parser.AlgExprAddSubContext) {
	l.BackToParentAlgebraicContext()
}

// EnterAlgExprMulDiv is called when production AlgExprMulDiv is entered.
func (l *RcalcParserListener) EnterAlgExprMulDiv(ctx *parser.AlgExprMulDivContext) {
	l.StartNewAlgebraicContext(&AlgebraicMulDivContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprMulDiv is called when production AlgExprMulDiv is exited.
func (l *RcalcParserListener) ExitAlgExprMulDiv(ctx *parser.AlgExprMulDivContext) {
	//fmt.Printf("ExitAlgExprMulDiv %s\n", ctx.GetText())
	l.BackToParentAlgebraicContext()
}

// EnterAlgExprFuncCall is called when production AlgExprFuncAtom is entered.
func (l *RcalcParserListener) EnterAlgExprFuncCall(ctx *parser.AlgExprFuncCallContext) {
	functionName := ctx.GetFunction_name().GetText()
	l.StartNewAlgebraicContext(&AlgebraicFunctionContext{
		AlgebraicExprContext: AlgebraicExprContext{reg: l.registry},
		functionName:         functionName,
	})
}

// ExitAlgExprFuncCall is called when production AlgExprFuncAtom is exited.
func (l *RcalcParserListener) ExitAlgExprFuncCall(ctx *parser.AlgExprFuncCallContext) {

	l.BackToParentAlgebraicContext()
}

// EnterAlgExprNumber is called when entering the AlgExprNumber production.
func (l *RcalcParserListener) EnterAlgExprNumber(ctx *parser.AlgExprNumberContext) {
	value, err := decimal.NewFromString(ctx.GetText())
	if err != nil {
		panic(fmt.Sprintf("Cannot parse number %s", ctx.GetText()))
	}
	l.StartNewAlgebraicContext(&AlgebraicNumberContext{
		AlgebraicExprContext: AlgebraicExprContext{reg: l.registry},
		value:                value,
	})
}

// ExitAlgExprNumber is called when exiting the AlgExprNumber production.
func (l *RcalcParserListener) ExitAlgExprNumber(ctx *parser.AlgExprNumberContext) {
	l.BackToParentAlgebraicContext()
}

// EnterAlgExprVariable is called when production AlgExprVariable is entered.
func (l *RcalcParserListener) EnterAlgExprVariable(ctx *parser.AlgExprVariableContext) {
	l.StartNewAlgebraicContext(&AlgebraicVariableNameContext{
		AlgebraicExprContext: AlgebraicExprContext{reg: l.registry},
		value:                ctx.GetText()})
}

// ExitAlgExprVariable is called when production AlgExprVariable is exited.
func (l *RcalcParserListener) ExitAlgExprVariable(ctx *parser.AlgExprVariableContext) {
	l.BackToParentAlgebraicContext()
}

// EnterAlgExprAddSignedAtom is called when production AlgExprAddSignedAtom is entered.
func (l *RcalcParserListener) EnterAlgExprAddSignedAtom(ctx *parser.AlgExprAddSignedAtomContext) {
	l.StartNewAlgebraicContext(&AlgebraicAtomContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprAddSignedAtom is called when production AlgExprAddSignedAtom is exited.
func (l *RcalcParserListener) ExitAlgExprAddSignedAtom(ctx *parser.AlgExprAddSignedAtomContext) {
	l.BackToParentAlgebraicContext()
}

// EnterAlgExprSubSignedAtom is called when production AlgExprSubSignedAtom is entered.
func (l *RcalcParserListener) EnterAlgExprSubSignedAtom(ctx *parser.AlgExprSubSignedAtomContext) {
	l.StartNewAlgebraicContext(&AlgebraicAtomContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprSubSignedAtom is called when production AlgExprSubSignedAtom is exited.
func (l *RcalcParserListener) ExitAlgExprSubSignedAtom(ctx *parser.AlgExprSubSignedAtomContext) {
	l.BackToParentAlgebraicContext()
}

// EnterAlgExprAtom is called when production AlgExprAtom is entered.
func (l *RcalcParserListener) EnterAlgExprAtom(ctx *parser.AlgExprAtomContext) {
	l.StartNewAlgebraicContext(&AlgebraicAtomContext{
		AlgebraicExprContext{reg: l.registry},
	})
}

// ExitAlgExprAtom is called when production AlgExprAtom is exited.
func (l *RcalcParserListener) ExitAlgExprAtom(ctx *parser.AlgExprAtomContext) {
	l.BackToParentAlgebraicContext()
}

// ExitInstrActionOrVarCall is called when exiting the InstrActionOrVarCall.
func (l *RcalcParserListener) ExitInstrActionOrVarCall(ctx *parser.InstrActionOrVarCallContext) {
	//fmt.Println("ExitInstrActionOrVarCall")
	action, err := parseAction(ctx.GetText(), l.registry)
	if err != nil {
		//ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
		l.AddAction(&VariableEvaluationActionDesc{varName: ctx.GetText()})
	} else {
		l.AddAction(action)
	}
}

// ExitDeclarationVariable is called when exiting the DeclarationVariable production.
func (l *RcalcParserListener) ExitDeclarationVariable(ctx *parser.DeclarationVariableContext) {
	//fmt.Println("ExitDeclarationVariable")
	action, err := parseAction(ctx.GetText(), l.registry)
	if err != nil {
		//ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
		l.AddVarName(ctx.GetText())
	} else {
		//TODO we should raise error here
		l.AddAction(action)
	}
}

// ExitInstrOp is called when production InstrOp is exited.
func (l *RcalcParserListener) ExitInstrOp(ctx *parser.InstrOpContext) {
	//fmt.Println("ExitInstrOp")
	action, err := parseAction(ctx.GetText(), l.registry)
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(action)
	}
}

func (l *RcalcParserListener) VisitTerminal(node antlr.TerminalNode) {
	l.TokenVisited(node.GetSymbol().GetTokenType())
	//fmt.Printf("VisitTerminal : #%s# / #%d#\n", node.GetSymbol().GetText(), node.GetSymbol().GetTokenType())
}

// EnterInstIfThenElse is called when entering the InstIfThenElse production.
func (l *RcalcParserListener) EnterInstIfThenElse(c *parser.InstIfThenElseContext) {
	l.StartNewContext(&IfThenElseContext{})
}

// ExitInstIfThenElse is called when entering the InstIfThenElse production.
func (l *RcalcParserListener) ExitInstIfThenElse(c *parser.InstIfThenElseContext) {
	l.BackToParentContext()
}

// EnterInstrStartNextLoop is called when production InstrStartNextLoop is entered.
func (l *RcalcParserListener) EnterInstrStartNextLoop(ctx *parser.InstrStartNextLoopContext) {
	loopContext := &StartEndLoopContext{}
	l.StartNewContext(loopContext)
}

// ExitInstrStartNextLoop is called when production InstrStartNextLoop is exited.
func (l *RcalcParserListener) ExitInstrStartNextLoop(ctx *parser.InstrStartNextLoopContext) {
	l.BackToParentContext()
}

// EnterInstrForNextLoop is called when exiting the InstrForNextLoop production.
func (l *RcalcParserListener) EnterInstrForNextLoop(c *parser.InstrForNextLoopContext) {
	loopContext := &ForNextLoopContext{}
	l.StartNewContext(loopContext)
}

// ExitInstrForNextLoop is called when exiting the InstrForNextLoop production.
func (l *RcalcParserListener) ExitInstrForNextLoop(c *parser.InstrForNextLoopContext) {
	l.BackToParentContext()
}

// EnterProgramDeclaration is called when entering the InstrProgramDeclaration production.
func (l *RcalcParserListener) EnterProgramDeclaration(c *parser.ProgramDeclarationContext) {
	l.StartNewContext(&ProgramContext{})
}

// ExitProgramDeclaration is called when entering the InstrProgramDeclaration production.
func (l *RcalcParserListener) ExitProgramDeclaration(c *parser.ProgramDeclarationContext) {
	l.BackToParentContext()
}

func (l *RcalcParserListener) EnterInstrLocalVarCreation(c *parser.InstrLocalVarCreationContext) {
	l.StartNewContext(&InstrLocalVarCreationContext{})
}

// ExitInstrLocalVarCreation is called when exiting the InstrLocalVarCreation production.
func (l *RcalcParserListener) ExitInstrLocalVarCreation(c *parser.InstrLocalVarCreationContext) {
	l.BackToParentContext()
}

/* Error Reporting */

type RcalcParserErrorListener struct {
	messages []string
}

func (el *RcalcParserErrorListener) HasErrors() bool {
	return len(el.messages) > 0

}

func (el *RcalcParserErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	antlrParser := recognizer.(antlr.Parser)
	stack := antlrParser.GetRuleInvocationStack(antlrParser.GetParserRuleContext())
	message := fmt.Sprintf("SyntaxError (%d, %d) : %s with stack %v", line, column, msg, stack)
	el.messages = append(el.messages, message)
}

func (el *RcalcParserErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	message := "ReportAmbiguity"
	el.messages = append(el.messages, message)
}

func (el *RcalcParserErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	message := "ReportAttemptingFullContext"
	el.messages = append(el.messages, message)
}

func (el *RcalcParserErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	message := "ReportContextSensitivity"
	el.messages = append(el.messages, message)
}

func ParseToActions(cmds string, lexerName string, registry *ActionRegistry) ([]Action, error) {

	is := antlr.NewInputStream(cmds)

	// Create the Lexer
	lexer := parser.NewRcalcLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewRcalcParser(stream)

	// Error Listener
	el := &RcalcParserErrorListener{}

	// Finally parse the expression (by walking the tree)
	var listener *RcalcParserListener = CreateRcalcParserListener(registry)
	//p.RemoveErrorListeners()
	p.AddErrorListener(el)
	parseResult := p.Start()
	if el.HasErrors() {
		return nil, fmt.Errorf("There are %d error(s):\n - %s", len(el.messages), strings.Join(el.messages, "\n - "))
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, parseResult)

	return listener.rootPc.GetItems(), nil
}
