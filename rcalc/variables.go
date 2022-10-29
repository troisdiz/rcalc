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

func (se *NumericVariable) asNumericVar() NumericVariable {
	return *se
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

func (se *BooleanVariable) asBooleanVar() BooleanVariable {
	return *se
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

type AlgebraicExpressionVariable struct {
	CommonVariable
	value string
}

func (se *AlgebraicExpressionVariable) String() string {
	return fmt.Sprintf("AlgebraicExpressionVariable(%v) type = %d", se.value, se.fType)
}

func (se *AlgebraicExpressionVariable) asIdentifierVar() AlgebraicExpressionVariable {
	return *se
}

func (se *AlgebraicExpressionVariable) display() string {
	return fmt.Sprintf("'%s'", se.value)
}

func CreateIdentifierVariable(value string) Variable {
	var result = AlgebraicExpressionVariable{
		CommonVariable: CommonVariable{fType: TYPE_ALG_EXPR},
		value:          value,
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
