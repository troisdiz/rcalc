package rcalc

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/shopspring/decimal"
	"strings"
	parser "troisdizaines.com/rcalc/rcalc/parser"
)

type Program struct {
	actions []Action
}

func (p *Program) Run(stack *Stack, system System) error {
	for _, action := range p.actions {
		err := action.Apply(system, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

type VariablePutOnStackActionDesc struct {
	value Variable
}

func (a *VariablePutOnStackActionDesc) NbArgs() int {
	return 0
}

func (a *VariablePutOnStackActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (a *VariablePutOnStackActionDesc) Apply(system System, stack *Stack) error {
	stack.Push(a.value)
	return nil
}

func (a *VariablePutOnStackActionDesc) OpCode() string {
	return "__hidden__" + "PutOnStack"
}

func (a *VariablePutOnStackActionDesc) String() string {
	return fmt.Sprintf("%s(%s)", a.OpCode(), a.value.String())
}

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

	actions []Action
}

func (pc *StartEndLoopContext) AddAction(action Action) {
	pc.actions = append(pc.actions, action)
}

func (pc *StartEndLoopContext) CreateFinalAction() Action {
	return &StartNextLoopActionDesc{actions: pc.actions}
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
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
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
	loopContext := &StartEndLoopContext{
		actions: nil,
	}
	l.StartNewContext(loopContext)
}

// ExitInstrStartNextLoop is called when production InstrStartNextLoop is exited.
func (l *RcalcParserListener) ExitInstrStartNextLoop(ctx *parser.InstrStartNextLoopContext) {
	l.BackToParentContext()
}

func ParseToActions(cmds string, lexerName string, registry *ActionRegistry) ([]Action, error) {

	is := antlr.NewInputStream(cmds)

	// Create the Lexer
	lexer := parser.NewRcalcLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewRcalcParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener *RcalcParserListener = CreateRcalcParserListener(registry)
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Start())
	return listener.rootPc.GetActions(), nil
}
