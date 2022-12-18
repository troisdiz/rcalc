package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
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

type AlgebraicExpressionNode interface {
	evaluate(variableReader VariableReader) *NumericVariable
}

type AlgExprMulDiv struct {
	items     []AlgebraicExpressionNode
	operators []int
}

var _ AlgebraicExpressionNode = (*AlgExprMulDiv)(nil)

func (a *AlgExprMulDiv) evaluate(variableReader VariableReader) *NumericVariable {
	result := decimal.NewFromInt(1)
	for _, it := range a.items {
		result = result.Mul(it.evaluate(variableReader).asNumericVar().value)
	}
	return CreateNumericVariable(result).asNumericVar()
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

func (a *AlgExprAddSub) evaluate(variableReader VariableReader) *NumericVariable {
	result := decimal.NewFromInt(0)
	for idx, it := range a.items {
		operator := OPERATOR_ADD
		if idx > 1 {
			operator = a.operators[idx-1]
		}
		switch operator {
		case OPERATOR_ADD:
			result = result.Add(it.evaluate(variableReader).asNumericVar().value)
		case OPERATOR_SUB:
			result = result.Sub(it.evaluate(variableReader).asNumericVar().value)
		}
	}
	return CreateNumericVariable(result).asNumericVar()
}

var _ AlgebraicExpressionNode = (*AlgExprAddSub)(nil)

type AlgExprNumber struct {
	value decimal.Decimal
}

var _ AlgebraicExpressionNode = (*AlgExprNumber)(nil)

func (a *AlgExprNumber) evaluate(variableReader VariableReader) *NumericVariable {
	return CreateNumericVariable(a.value).asNumericVar()
}

type AlgExprVariable struct {
	value string
}

var _ AlgebraicExpressionNode = (*AlgExprVariable)(nil)

func (aev *AlgExprVariable) evaluate(variableReader VariableReader) *NumericVariable {
	variableValue, err := variableReader.GetVariableValue(aev.value)
	if err != nil {
		// TODO Error system for such case
		return nil
	}
	if variableValue.getType() == TYPE_NUMERIC {
		numericVar := variableValue.asNumericVar()
		return numericVar
	} else {
		// TODO Error system for such case
		return nil
	}
}

type AlgExprSignedElt struct {
	items    AlgebraicExpressionNode
	operator int
}

var _ AlgebraicExpressionNode = (*AlgExprSignedElt)(nil)

func (a *AlgExprSignedElt) evaluate(variableReader VariableReader) *NumericVariable {
	result := a.items.evaluate(variableReader)
	if a.operator == OPERATOR_SUB {
		result.value = result.value.Neg()
	}
	return result
}

type AlgExprFunctionElt struct {
	//function     interface{}
	functionName string
	arguments    []AlgebraicExpressionNode
}

var _ AlgebraicExpressionNode = (*AlgExprFunctionElt)(nil)

func (a AlgExprFunctionElt) evaluate(variableReader VariableReader) *NumericVariable {

	if a.functionName == "cos" {
		if len(a.arguments) == 1 {
			argValue := a.arguments[0].evaluate(variableReader)
			result := argValue.value.Cos()
			return CreateNumericVariable(result).asNumericVar()
		}
	}

	//TODO implement me
	panic("implement me")
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
	return fmt.Sprintf(" << %s >>", strings.Join(actionStr, " "))
}

func (p *ProgramVariable) asProgramVar() *ProgramVariable {
	return p
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
	default:
		return nil, fmt.Errorf("unknown variable type")
	}
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
	default:
		return nil, fmt.Errorf("marshalling of programs not implemented yet")
	}

}
