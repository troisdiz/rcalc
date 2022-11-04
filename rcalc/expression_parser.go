package rcalc

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/shopspring/decimal"
	"strings"
	parser "troisdizaines.com/rcalc/rcalc/parser"
)

func parserNumber(txt string) (Variable, error) {
	number, err := decimal.NewFromString(txt)
	if err != nil {
		return nil, err
	}
	return CreateNumericVariable(number), nil
}

func parseIdentifier(txt string) (Variable, error) {
	l := len(txt)
	return CreateIdentifierVariable(txt[1 : l-1]), nil
}

func parseAction(txt string, registry *ActionRegistry) (Action, error) {
	lowerTxt := strings.ToLower(txt)
	if registry.ContainsOpCode(lowerTxt) {
		return registry.GetAction(txt), nil
	} else {
		return nil, fmt.Errorf("unknown action")
	}
}

type ParseContext[T any] interface {
	GetParent() ParseContext[T]
	SetParent(ctx ParseContext[T])

	AddAction(action T)
	AddIdentifier(id string)

	BackFromChild(child ParseContext[T])
	CreateFinalAction() T

	TokenVisited(token int)
}

type BaseParseContext[T any] struct {
	parent         ParseContext[T]
	actions        []T
	idDeclarations []string
}

var _ ParseContext[string] = (*BaseParseContext[string])(nil)

func (g *BaseParseContext[T]) GetParent() ParseContext[T] {
	return g.parent
}

func (g *BaseParseContext[T]) SetParent(ctx ParseContext[T]) {
	g.parent = ctx
}

func (g *BaseParseContext[T]) AddAction(action T) {
	g.actions = append(g.actions, action)
}

func (g *BaseParseContext[T]) AddIdentifier(id string) {
	g.idDeclarations = append(g.idDeclarations, id)
}

func (g *BaseParseContext[T]) BackFromChild(child ParseContext[T]) {
	g.AddAction(child.CreateFinalAction())
}

func (g *BaseParseContext[T]) CreateFinalAction() T {
	panic("CreateFinalAction must be implemented by sub structures")
}

func (g *BaseParseContext[T]) GetActions() []T {
	return g.actions
}

func (g *BaseParseContext[T]) TokenVisited(token int) {}

type IfThenElseContext struct {
	BaseParseContext[Action] // to avoid reimplementing the interface

	actions       [][]Action
	currentAction int
}

func (i *IfThenElseContext) AddAction(action Action) {
	i.actions[i.currentAction] = append(i.actions[i.currentAction], action)
}

func (i *IfThenElseContext) TokenVisited(token int) {
	switch token {
	case parser.RcalcLexerKW_IF:
		i.actions = make([][]Action, 3)
		i.currentAction = 0
	case parser.RcalcLexerKW_THEN:
		i.currentAction = 1
	case parser.RcalcLexerKW_ELSE:
		i.currentAction = 2
	}
}

func (i *IfThenElseContext) CreateFinalAction() Action {
	return &IfThenElseActionDesc{
		ifActions:   i.actions[0],
		thenActions: i.actions[1],
		elseActions: i.actions[2],
	}
}

type StartEndLoopContext struct {
	BaseParseContext[Action]
}

func (pc *StartEndLoopContext) CreateFinalAction() Action {
	return &StartNextLoopActionDesc{actions: pc.BaseParseContext.actions}
}

type ForNextLoopContext struct {
	BaseParseContext[Action]
}

func (pc *ForNextLoopContext) CreateFinalAction() Action {
	return &ForNextLoopActionDesc{
		varName: pc.BaseParseContext.idDeclarations[0],
		actions: pc.BaseParseContext.actions,
	}
}

type ProgramContext struct {
	BaseParseContext[Action]
}

func (pc *ProgramContext) CreateFinalAction() Action {
	progVar := CreateProgramVariable(pc.actions)
	return &VariablePutOnStackActionDesc{value: progVar}
}

type RcalcParserListener struct {
	*parser.BaseRcalcListener

	registry *ActionRegistry

	rootPc    *BaseParseContext[Action]
	currentPc ParseContext[Action]
}

func CreateRcalcParserListener(registry *ActionRegistry) *RcalcParserListener {
	rootPc := &BaseParseContext[Action]{
		parent:  nil,
		actions: nil,
	}
	return &RcalcParserListener{
		registry:  registry,
		rootPc:    rootPc,
		currentPc: rootPc,
	}
}

func (l *RcalcParserListener) AddAction(action Action) {
	l.currentPc.AddAction(action)
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

func (l *RcalcParserListener) TokenVisited(token int) {
	l.currentPc.TokenVisited(token)
}

// ExitVariableNumber is called when production InstrNumber is exited.
func (l *RcalcParserListener) ExitVariableNumber(ctx *parser.VariableNumberContext) {
	fmt.Printf("ExitInstrNumber: %s\n", ctx.GetText())
	number, err := parserNumber(ctx.GetText())
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(&VariablePutOnStackActionDesc{number})
	}

}

type AlgebraicVariableContext struct {
	BaseParseContext[Action] // to avoid reimplementing the interface
}

func (ac *AlgebraicVariableContext) CreateFinalAction() Action {
	return nil
}

func (ac *AlgebraicVariableContext) TokenVisited(token int) {

}

// EnterVariableAlgebraicExpression is called when production VariableAlgebraicExpression is entered.
func (l *RcalcParserListener) EnterVariableAlgebraicExpression(ctx *parser.VariableAlgebraicExpressionContext) {
	fmt.Println("EnterVariableAlgebraicExpression")
	l.StartNewContext(&AlgebraicVariableContext{})
}

// ExitVariableAlgebraicExpression is called when production VariableAlgebraicExpression is exited.
func (l *RcalcParserListener) ExitVariableAlgebraicExpression(ctx *parser.VariableAlgebraicExpressionContext) {
	fmt.Println("ExitVariableAlgebraicExpression")
	identifier, err := parseIdentifier(ctx.GetText())
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(&VariablePutOnStackActionDesc{value: identifier})
	}
	l.BackToParentContext()
}

// ExitInstrActionOrVarCall is called when exiting the InstrActionOrVarCall.
func (l *RcalcParserListener) ExitInstrActionOrVarCall(ctx *parser.InstrActionOrVarCallContext) {
	fmt.Println("ExitInstrActionOrVarCall")
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
	fmt.Println("ExitInstrActionOrVarCall")
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
	fmt.Println("ExitInstrOp")
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

// EnterInstrProgramDeclaration is called when entering the InstrProgramDeclaration production.
func (l *RcalcParserListener) EnterInstrProgramDeclaration(c *parser.InstrProgramDeclarationContext) {
	l.StartNewContext(&ProgramContext{})
}

// ExitInstrProgramDeclaration is called when entering the InstrProgramDeclaration production.
func (l *RcalcParserListener) ExitInstrProgramDeclaration(c *parser.InstrProgramDeclarationContext) {
	//programContext := l.currentPc
	l.BackToParentContext()
	//l.AddAction(programContext.CreateFinalAction())
}

// ExitAlgExprRoot is called when production AlgExprRoot is exited.
func (l *RcalcParserListener) ExitAlgExprRoot(ctx *parser.AlgExprRootContext) {
	op_type := ctx.GetOp_type()
	if op_type != nil {
		fmt.Printf("ExitAlgExprRoot => %s\n", op_type.GetText())
	}
	fmt.Printf("ExitAlgExprRoot full text => %s\n", ctx.GetText())
}

/* Error Reporting */

type RcalcParserErrorListener struct {
	messages []string
}

func (el *RcalcParserErrorListener) HasErrors() bool {
	return len(el.messages) > 0

}

func (el *RcalcParserErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	message := fmt.Sprintf("SyntaxError (%d, %d) : %s", line, column, msg)
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
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Start())
	if el.HasErrors() {

		return nil, fmt.Errorf("There are %d error(s):\n - %s", len(el.messages), strings.Join(el.messages, "\n - "))
	}
	return listener.rootPc.GetActions(), nil
}
