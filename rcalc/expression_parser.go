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
	var lexItems []LexItem

	for lextItem := lexer.NextItem(); lextItem.typ != lexItemEOF; lextItem = lexer.NextItem() {
		lexItems = append(lexItems, lextItem)
	}
	for idx := 0; idx < len(lexItems); idx++ {
		lexItem := lexItems[idx]
		switch lexItem.typ {
		case lexItemNumber:
			variable, err := parserNumber(lexItem)
			if err != nil {
				return nil, err
			}
			result = append(result, &VariablePutOnStackActionDesc{value: variable})
		case lexItemAction:
			action, err := parseAction(lexItem, registry)
			if err != nil {
				result = append(result, action)
			}
		case lexItemIdentifier:
			variable, err := parseIdentifier(lexItem)
			if err != nil {
				return nil, err
			}
			result = append(result, &VariablePutOnStackActionDesc{value: variable})
		case lexItemOpenCurlyBrace:
			list, closeIdx, err := parseList(lexItems[idx:], registry)
			if err != nil {
				return nil, err
			}
			result = append(result, &VariablePutOnStackActionDesc{value: list})
			idx = closeIdx
		default:
			fmt.Printf("Ignore %v for now\n", lexItem)
		}
	}
	return result, nil
}
func parseList(items []LexItem, registry *ActionRegistry) (Variable, int, error) {

	lastIdemIdx := len(items) - 1
	endIdx := -1
	for i := 0; i <= lastIdemIdx; i++ {
		if items[i].typ == lexItemCloseCurlyBrace {
			endIdx = i
			break
		} else if i == lastIdemIdx {
			return nil, endIdx, fmt.Errorf("list started is not closed")
		}
	}
	// end has been found, lets create the list
	variables := make([]Variable, endIdx-1)
	for idx, item := range items[0:endIdx] {
		variable, err := parseVariable(item, registry)
		if err != nil {
			return nil, -1, err
		}
		variables[idx] = variable
	}
	// TODO return
}

func parseVariable(lexItem LexItem, registry *ActionRegistry) (Variable, error) {
	switch lexItem.typ {
	case lexItemNumber:
		number, err := parserNumber(lexItem)
		if err != nil {
			return nil, err
		}
		return number, nil
	case lexItemIdentifier:
		identifier, err := parseIdentifier(lexItem)
		if err != nil {
			return nil, err
		}
		return identifier, nil
	default:
		return nil, fmt.Errorf("%s cannot be parsed")
	}
}

func parserNumber(lexItem LexItem) (Variable, error) {
	number, err := decimal.NewFromString(lexItem.val)
	if err != nil {
		return nil, err
	}
	return CreateNumericVariable(number), nil
}

func parseIdentifier(item LexItem) (Variable, error) {
	l := len(item.val)
	return CreateIdentifierVariable(item.val[1 : l-1]), nil
}

func parseAction(lexItem LexItem, registry *ActionRegistry) (Action, error) {
	if registry.ContainsOpCode(lexItem.val) {
		return registry.GetAction(lexItem.val), nil
	} else {
		return nil, fmt.Errorf("unknow action")
	}
}
