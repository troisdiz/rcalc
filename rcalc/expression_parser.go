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

func ParseToActions(cmds string, lexerName string, registry *ActionRegistry) ([]Action, error) {
	lexer := Lex(lexerName, cmds)
	var result []Action
	for lextItem := lexer.NextItem(); lextItem.typ != lexItemEOF; lextItem = lexer.NextItem() {
		switch lextItem.typ {
		case lexItemNumber:
			number, err := decimal.NewFromString(lextItem.val)
			if err != nil {
				return nil, err
			}
			result = append(result, &VariablePutOnStackActionDesc{value: CreateNumericVariable(number)})
		case lexItemAction:
			if registry.ContainsOpCode(lextItem.val) {
				result = append(result, registry.GetAction(lextItem.val))
			}
		case lexItemIdentifier:
			l := len(lextItem.val)
			variable := CreateIdentifierVariable(lextItem.val[1 : l-1])
			result = append(result, &VariablePutOnStackActionDesc{value: variable})
		default:
			fmt.Printf("Ignore %v for now\n", lextItem)
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

type RcalcParserListener struct {
	*parser.BaseRcalcListener

	registry *ActionRegistry

	actions []Action
}

func (l *RcalcParserListener) AddAction(action Action) {
	l.actions = append(l.actions, action)
}

// ExitInstrNumber is called when production InstrNumber is exited.
func (l *RcalcParserListener) ExitInstrNumber(ctx *parser.InstrNumberContext) {
	fmt.Printf("ExitInstrNumber: %s\n", ctx.GetText())
	number, err := parserNumber(ctx.GetText())
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(&VariablePutOnStackActionDesc{number})
	}

}

// ExitInstrIdentifier is called when production identifier is exited.
func (l *RcalcParserListener) ExitIdentifier(ctx *parser.IdentifierContext) {
	fmt.Println("ExitInstrIdentifier")
	identifier, err := parseIdentifier(ctx.GetText())
	if err != nil {
		ctx.AddErrorNode(ctx.GetParser().GetCurrentToken())
	} else {
		l.AddAction(&VariablePutOnStackActionDesc{value: identifier})
	}
}

// ExitInstrActionOrVarCall is called when exiting the InstrActionOrVarCall.
func (l *RcalcParserListener) ExitAction_or_var_call(ctx *parser.Action_or_var_callContext) {
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

func ParseToActions2(cmds string, lexerName string, registry *ActionRegistry) ([]Action, error) {

	is := antlr.NewInputStream(cmds)

	// Create the Lexer
	lexer := parser.NewRcalcLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewRcalcParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener RcalcParserListener = RcalcParserListener{registry: registry}
	antlr.ParseTreeWalkerDefault.Walk(&listener, p.Start())
	return listener.actions, nil
}
