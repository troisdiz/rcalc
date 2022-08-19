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

type ParseContext interface {
	GetParent() ParseContext
	SetParent(ctx ParseContext)

	AddAction(action Action)
	AddIdentifier(id string)

	BackFromChild(child ParseContext)
	CreateFinalAction() Action
}

type BaseParseContext struct {
	parent         ParseContext
	actions        []Action
	idDeclarations []string
}

func (pc *BaseParseContext) CreateFinalAction() Action {
	panic("CreateFinalAction must be implemented by sub structures")
}

func (pc *BaseParseContext) GetParent() ParseContext {
	return pc.parent
}

func (pc *BaseParseContext) SetParent(ctx ParseContext) {
	pc.parent = ctx
}

func (pc *BaseParseContext) BackFromChild(child ParseContext) {
	pc.AddAction(child.CreateFinalAction())
}

func (pc *BaseParseContext) AddAction(action Action) {
	pc.actions = append(pc.actions, action)
}

func (pc *BaseParseContext) AddIdentifier(id string) {
	pc.idDeclarations = append(pc.idDeclarations, id)
}

func (pc *BaseParseContext) GetActions() []Action {
	return pc.actions
}

type StartEndLoopContext struct {
	BaseParseContext
}

func (pc *StartEndLoopContext) CreateFinalAction() Action {
	return &StartNextLoopActionDesc{actions: pc.BaseParseContext.actions}
}

type ForNextLoopContext struct {
	BaseParseContext
}

func (pc *ForNextLoopContext) CreateFinalAction() Action {
	return &ForNextLoopActionDesc{
		varName: pc.BaseParseContext.idDeclarations[0],
		actions: pc.BaseParseContext.actions,
	}
}

type RcalcParserListener struct {
	*parser.BaseRcalcListener

	registry *ActionRegistry

	rootPc    *BaseParseContext
	currentPc ParseContext
}

func CreateRcalcParserListener(registry *ActionRegistry) *RcalcParserListener {
	rootPc := &BaseParseContext{
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

func (l *RcalcParserListener) StartNewContext(ctx ParseContext) {
	ctx.SetParent(l.currentPc)
	l.currentPc = ctx
}

func (l *RcalcParserListener) BackToParentContext() {
	oldCurrent := l.currentPc
	l.currentPc = l.currentPc.GetParent()
	l.currentPc.BackFromChild(oldCurrent)
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

// ExitVariableIdentifier is called when production identifier is exited.
func (l *RcalcParserListener) ExitVariableIdentifier(ctx *parser.VariableIdentifierContext) {
	fmt.Println("ExitInstrIdentifier")
	identifier, err := parseIdentifier(ctx.GetText())
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(&VariablePutOnStackActionDesc{value: identifier})
	}
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
		return nil, fmt.Errorf("There are %d error(s)", len(el.messages))
	}
	return listener.rootPc.GetActions(), nil
}
