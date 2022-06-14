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

func ParseToActions(lexer *Lexer, registry *ActionRegistry) ([]Action, error) {
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
