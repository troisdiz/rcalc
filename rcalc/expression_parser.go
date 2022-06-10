package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
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

type DecimalPutOnStackActionDesc struct {
	number decimal.Decimal
}

func (a *DecimalPutOnStackActionDesc) NbArgs() int {
	return 0
}

func (a *DecimalPutOnStackActionDesc) CheckTypes(elts ...StackElt) (bool, error) {
	return true, nil
}

func (a *DecimalPutOnStackActionDesc) Apply(system System, stack *Stack) error {
	stack.Push(CreateNumericStackElt(a.number))
	return nil
}

func (a *DecimalPutOnStackActionDesc) OpCode() string {
	return "__hidden__"
}

func ParseToActions(lexer *Lexer, registry *ActionRegistry, ) ([]Action, error) {
	var result []Action
	for lextItem := lexer.NextItem(); lextItem.typ != lexItemEOF; lextItem = lexer.NextItem() {
		switch lextItem.typ {
		case lexItemNumber:
			number, err := decimal.NewFromString(lextItem.val)
			if err != nil {
				return nil, err
			}
			result = append(result, &DecimalPutOnStackActionDesc{number: number})
		case lexItemAction:
			if registry.ContainsOpCode(lextItem.val) {
				result = append(result, registry.GetAction(lextItem.val))
			}
		default:
			fmt.Printf("Ignore %v for now\n", lextItem)
		}
	}
	return result, nil
}
