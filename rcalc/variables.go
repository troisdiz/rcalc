package rcalc

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

type NumericVariable struct {
	CommonVariable
	value decimal.Decimal
}

func (se *NumericVariable) String() string {
	return fmt.Sprintf("NumericVariable(%v)", se.value)
}

func (se *NumericVariable) asNumericVar() *NumericVariable {
	return se
}

func (se *NumericVariable) display() string {
	return se.value.String()
}

func CreateNumericVariable(value decimal.Decimal) Variable {
	var result = NumericVariable{
		CommonVariable: CommonVariable{fType: TYPE_NUMERIC},
		value:          value,
	}
	return &result
}

func CreateNumericVariableFromInt(value int) Variable {
	var result = NumericVariable{
		CommonVariable: CommonVariable{fType: TYPE_NUMERIC},
		value:          decimal.NewFromInt(int64(value)),
	}
	return &result
}

func (se *NumericVariable) Equals(other *NumericVariable) bool {
	if other != nil {
		return se.value.Equals(other.value)
	}
	return false
}

type BooleanVariable struct {
	CommonVariable
	value bool
}

func (se *BooleanVariable) String() string {
	return fmt.Sprintf("BooleanVariable(%v) type = %d", se.value, se.fType)
}

func (se *BooleanVariable) asBooleanVar() *BooleanVariable {
	return se
}

func (se *BooleanVariable) display() string {
	return fmt.Sprintf("%t", se.value)
}

func CreateBooleanVariable(value bool) Variable {
	var result = BooleanVariable{
		CommonVariable: CommonVariable{fType: TYPE_BOOL},
		value:          value,
	}
	return &result
}

type ListVariable struct {
	CommonVariable
	items []Variable
}

var _ Variable = (*ListVariable)(nil)

func CreateListVariable(items []Variable) *ListVariable {
	return &ListVariable{
		CommonVariable: CommonVariable{fType: TYPE_LIST},
		items:          items,
	}
}

func (l *ListVariable) display() string {
	displayedItems := make([]string, len(l.items))
	for i, item := range l.items {
		displayedItems[i] = item.display()
	}
	return fmt.Sprintf("{ %s }", strings.Join(displayedItems, " "))
}

func (l *ListVariable) asListVar() *ListVariable {
	return l
}

func (l *ListVariable) Size() int {
	return len(l.items)
}

type AlgebraicExpressionNode interface {
	Evaluate(variableReader VariableReader) (*NumericVariable, error)
}

type AlgExprMulDiv struct {
	items     []AlgebraicExpressionNode
	operators []int
}

var _ AlgebraicExpressionNode = (*AlgExprMulDiv)(nil)

func (a *AlgExprMulDiv) Evaluate(variableReader VariableReader) (*NumericVariable, error) {

	if len(a.items) == 1 {
		return a.items[0].Evaluate(variableReader)
	}

	result := decimal.NewFromInt(1)
	for _, it := range a.items {
		variable, err := it.Evaluate(variableReader)
		if err != nil {
			return nil, err
		}
		numericVar := variable.asNumericVar()
		subExprValue := numericVar.value
		result = result.Mul(subExprValue)
	}
	return CreateNumericVariable(result).asNumericVar(), nil
}

type AlgExprPow struct {
	items []AlgebraicExpressionNode
}

var _ AlgebraicExpressionNode = (*AlgExprPow)(nil)

func (a *AlgExprPow) Evaluate(variableReader VariableReader) (*NumericVariable, error) {

	if len(a.items) == 1 {
		return a.items[0].Evaluate(variableReader)
	}

	fmt.Printf("Pow items(%d):\n", len(a.items))
	for idx, it := range a.items {
		fmt.Printf("  - item[%d] = %v\n", idx, it)
	}

	result := decimal.NewFromInt(1)
	for i := len(a.items) - 1; i >= 0; i-- {
		fmt.Printf("i = %d\n", i)
		variable, err := a.items[i].Evaluate(variableReader)
		if err != nil {
			return nil, err
		}
		numericVar := variable.asNumericVar()
		if result.Equal(decimal.NewFromInt(1)) {
			result = numericVar.value
		} else {
			result = numericVar.value.Pow(result)
		}
	}
	return CreateNumericVariable(result).asNumericVar(), nil
}

// type AddSubOp int
const (
	OPERATOR_ADD = 0
	OPERATOR_SUB = 1
)

type AlgExprAddSub struct {
	items     []AlgebraicExpressionNode
	operators []int
}

func (a *AlgExprAddSub) Evaluate(variableReader VariableReader) (*NumericVariable, error) {

	if len(a.items) == 1 {
		return a.items[0].Evaluate(variableReader)
	}

	result := decimal.NewFromInt(0)
	for idx, it := range a.items {
		operator := OPERATOR_ADD
		if idx > 1 {
			operator = a.operators[idx-1]
		}
		switch operator {
		case OPERATOR_ADD:
			variable, err := it.Evaluate(variableReader)
			if err != nil {
				return nil, err
			}
			result = result.Add(variable.asNumericVar().value)
		case OPERATOR_SUB:
			variable, err := it.Evaluate(variableReader)
			if err != nil {
				return nil, err
			}
			result = result.Sub(variable.asNumericVar().value)
		}
	}
	return CreateNumericVariable(result).asNumericVar(), nil
}

var _ AlgebraicExpressionNode = (*AlgExprAddSub)(nil)

type AlgExprNumber struct {
	value decimal.Decimal
}

var _ AlgebraicExpressionNode = (*AlgExprNumber)(nil)

func (a *AlgExprNumber) Evaluate(variableReader VariableReader) (*NumericVariable, error) {
	return CreateNumericVariable(a.value).asNumericVar(), nil
}

type AlgExprVariable struct {
	value string
}

var _ AlgebraicExpressionNode = (*AlgExprVariable)(nil)

func (aev *AlgExprVariable) Evaluate(variableReader VariableReader) (*NumericVariable, error) {
	varName := aev.value
	variableValue, err := variableReader.GetVariableValue(varName)
	if err != nil {
		return nil, fmt.Errorf("cannot find variable %s", varName)
	}
	if variableValue.getType() == TYPE_NUMERIC {
		numericVar := variableValue.asNumericVar()
		return numericVar, nil
	} else {
		return nil, fmt.Errorf("variable %s is not of numeric type", varName)
	}
}

type AlgExprSignedElt struct {
	items    AlgebraicExpressionNode
	operator int
}

var _ AlgebraicExpressionNode = (*AlgExprSignedElt)(nil)

func (a *AlgExprSignedElt) Evaluate(variableReader VariableReader) (*NumericVariable, error) {
	result, _ := a.items.Evaluate(variableReader)
	if a.operator == OPERATOR_SUB {
		result.value = result.value.Neg()
	}
	return result, nil
}

type AlgExprFunctionElt struct {
	//function     interface{}
	functionName string
	fn           AlgebraicFn
	arguments    []AlgebraicExpressionNode
}

var _ AlgebraicExpressionNode = (*AlgExprFunctionElt)(nil)

func (a AlgExprFunctionElt) Evaluate(variableReader VariableReader) (*NumericVariable, error) {

	var args []decimal.Decimal = make([]decimal.Decimal, len(a.arguments))
	for idx, numVarArg := range a.arguments {
		value, err := numVarArg.Evaluate(variableReader)
		if err != nil {
			return nil, err
		} else {
			args[idx] = value.value
		}
	}

	return CreateNumericVariable(a.fn(args...)).asNumericVar(), nil
}

type AlgebraicExpressionVariable struct {
	CommonVariable
	value    string
	rootNode AlgebraicExpressionNode
}

func (se *AlgebraicExpressionVariable) String() string {
	return fmt.Sprintf("AlgebraicExpressionVariable(%v) type = %d", se.value, se.fType)
}

func (se *AlgebraicExpressionVariable) asIdentifierVar() *AlgebraicExpressionVariable {
	return se
}

func (se *AlgebraicExpressionVariable) display() string {
	return fmt.Sprintf("'%s'", se.value)
}

func CreateAlgebraicExpressionVariable(value string, algExprNode AlgebraicExpressionNode) Variable {
	var result = AlgebraicExpressionVariable{
		CommonVariable: CommonVariable{fType: TYPE_ALG_EXPR},
		value:          value,
		rootNode:       algExprNode,
	}
	return &result
}

type ProgramVariable struct {
	CommonVariable
	actions []Action
}

func (p *ProgramVariable) display() string {

	actionStr := []string{}
	for _, action := range p.actions {
		actionStr = append(actionStr, action.Display())
	}
	return fmt.Sprintf("<< %s >>", strings.Join(actionStr, " "))
}

func (p *ProgramVariable) asProgramVar() *ProgramVariable {
	return p
}

func (p *ProgramVariable) String() string {
	var actionStrings []string
	for _, action := range p.actions {
		actionStrings = append(actionStrings, fmt.Sprintf("%v", action))
	}
	return fmt.Sprintf("[ProgramVariable]\n    %s", strings.Join(actionStrings, "\n    "))
}

func CreateProgramVariable(actions []Action) *ProgramVariable {
	return &ProgramVariable{
		CommonVariable: CommonVariable{
			fType: TYPE_PROGRAM,
		},
		actions: actions,
	}
}

func CreateProtoFromProgram(prg *ProgramVariable) (*protostack.ProgramVariable, error) {
	protoProgram := &protostack.ProgramVariable{}
	for _, action := range prg.actions {
		protoAction, err := action.MarshallFunc()(nil, action)
		if err != nil {
			return nil, err
		}
		protoProgram.Actions = append(protoProgram.Actions, protoAction)
	}
	return protoProgram, nil
}

func CreateProgramVariableFromProto(
	reg *ActionRegistry,
	protoProgramVariable *protostack.ProgramVariable) (*ProgramVariable, error) {

	var actions []Action
	for _, protoAction := range protoProgramVariable.GetActions() {
		action, err := reg.CreateActionFromProto(protoAction)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return CreateProgramVariable(actions), nil
}

func CreateProtoFromAlgExpr(algExpr *AlgebraicExpressionVariable) (*protostack.AlgebraicExpressionVariable, error) {
	protoAlgExpr := &protostack.AlgebraicExpressionVariable{
		FullText: algExpr.value,
	}
	return protoAlgExpr, nil
}

func CreateListFromProto(reg *ActionRegistry, protoListVariable *protostack.ListVariable) (Variable, error) {

	var items []Variable
	for _, protoVar := range protoListVariable.GetItems() {
		action, err := CreateVariableFromProto(reg, protoVar)
		if err != nil {
			return nil, err
		}
		items = append(items, action)
	}
	return CreateListVariable(items), nil
}

func CreateVariableFromProto(reg *ActionRegistry, protoVariable *protostack.Variable) (Variable, error) {
	switch protoVariable.GetType() {
	case protostack.VariableType_NUMBER:
		protoNumber := protoVariable.GetNumber()
		decimalNumber := decimal.NewFromInt(0)
		err := decimalNumber.UnmarshalBinary(protoNumber.GetValue())
		if err != nil {
			return nil, err
		}
		return CreateNumericVariable(decimalNumber), nil
	case protostack.VariableType_BOOLEAN:
		return CreateBooleanVariable(protoVariable.GetBool().GetValue()), nil
	case protostack.VariableType_PROGRAM:
		return CreateProgramVariableFromProto(reg, protoVariable.GetProgram())
	case protostack.VariableType_ALGEBRAIC_EXPRESSION:
		return CreateAlgebraicExpressionVariableFromProto(reg, protoVariable.GetAlgExpr())
	case protostack.VariableType_LIST:
		return CreateListFromProto(reg, protoVariable.GetList())
	default:
		return nil, fmt.Errorf("unknown variable type")
	}
}

func CreateAlgebraicExpressionVariableFromProto(
	reg *ActionRegistry,
	protoAlgExpr *protostack.AlgebraicExpressionVariable) (*AlgebraicExpressionVariable, error) {
	actions, err := ParseToActions(fmt.Sprintf("'%s'", protoAlgExpr.FullText), "", reg)
	return actions[0].(*VariablePutOnStackActionDesc).value.(*AlgebraicExpressionVariable), err
}

func CreateProtoFromVariable(variable Variable) (*protostack.Variable, error) {
	switch variable.getType() {
	case TYPE_NUMERIC:
		binaryNumber, err := variable.asNumericVar().value.MarshalBinary()
		if err != nil {
			return nil, err
		}
		protoNumVar := &protostack.NumberVariable{Value: binaryNumber}
		protoVar := &protostack.Variable{
			Type:    protostack.VariableType_NUMBER,
			RealVar: &protostack.Variable_Number{Number: protoNumVar},
		}
		return protoVar, nil
	case TYPE_BOOL:
		protoBoolVar := &protostack.BooleanVariable{Value: variable.asBooleanVar().value}
		return &protostack.Variable{
			Type:    protostack.VariableType_BOOLEAN,
			RealVar: &protostack.Variable_Bool{Bool: protoBoolVar},
		}, nil
	case TYPE_PROGRAM:
		protoProgramVar, err := CreateProtoFromProgram(variable.asProgramVar())
		if err != nil {
			return nil, err
		}
		return &protostack.Variable{
				Type:    protostack.VariableType_PROGRAM,
				RealVar: &protostack.Variable_Program{Program: protoProgramVar}},
			nil
	case TYPE_ALG_EXPR:
		algExprVar := variable.(*AlgebraicExpressionVariable)

		protoAlgExpr, err := CreateProtoFromAlgExpr(algExprVar)
		if err != nil {
			return nil, err
		}
		return &protostack.Variable{
			Type:    protostack.VariableType_ALGEBRAIC_EXPRESSION,
			RealVar: &protostack.Variable_AlgExpr{AlgExpr: protoAlgExpr},
		}, nil
	case TYPE_LIST:
		listVar := variable.(*ListVariable)
		protoItems := make([]*protostack.Variable, len(listVar.items))
		for i := 0; i < len(listVar.items); i++ {
			protoItem, err := CreateProtoFromVariable(listVar.items[i])
			if err != nil {
				return nil, err
			}
			protoItems[i] = protoItem
		}
		protoListVar := &protostack.ListVariable{Items: protoItems}
		return &protostack.Variable{
			Type:    protostack.VariableType_LIST,
			RealVar: &protostack.Variable_List{List: protoListVar},
		}, nil
	default:
		return nil, fmt.Errorf("marshalling of variables of type %d is not implemented yet", variable.getType())
	}

}
