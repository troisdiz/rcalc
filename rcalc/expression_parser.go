package rcalc

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/shopspring/decimal"
	parser "troisdizaines.com/rcalc/rcalc/parser/grammar"
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

type RcalcParserListener struct {
	*parser.BaseRcalcListener
}

func ParseToActions2(cmds string, lexerName string, registry *ActionRegistry) ([]Action, error) {

	is := antlr.NewInputStream(cmds)

	// Create the Lexer
	lexer := parser.NewRcalcLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewRcalcParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener RcalcParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, p.Start())
	return nil, fmt.Errorf("to be implemented")
}
