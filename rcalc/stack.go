package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
	"os"
	"strings"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

type Type int

const (
	TYPE_GENERIC    Type = 0
	TYPE_NUMERIC    Type = 1
	TYPE_BOOL       Type = 2
	TYPE_STR        Type = 3
	TYPE_IDENTIFIER Type = 4
	TYPE_PROGRAM    Type = 5
	// TYPE_LIST       Type = 6
	// TYPE_VECTOR     Type = 7
)

type Variable interface {
	getType() Type
	asNumericVar() NumericVariable
	asBooleanVar() BooleanVariable
	asIdentifierVar() IdentifierVariable
	asProgramVar() *ProgramVariable
	display() string
	String() string
}

type CommonVariable struct {
	fType Type
}

func (se *CommonVariable) getType() Type {
	return se.fType
}

func (se *CommonVariable) asNumericVar() NumericVariable {
	panic("This is not a Numeric variable")
}

func (se *CommonVariable) asBooleanVar() BooleanVariable {
	panic("This is not a Boolean variable")
}

func (se *CommonVariable) asIdentifierVar() IdentifierVariable {
	panic("This is not an Identifier variable")
}

func (se *CommonVariable) asProgramVar() *ProgramVariable {
	panic("This is not a Program variable")
}

func (se *CommonVariable) String() string {
	return fmt.Sprintf("[CommonVariable] t=%d", se.fType)
}

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

type IdentifierVariable struct {
	CommonVariable
	value string
}

func (se *IdentifierVariable) String() string {
	return fmt.Sprintf("IdentifierVariable(%v) type = %d", se.value, se.fType)
}

func (se *IdentifierVariable) asIdentifierVar() IdentifierVariable {
	return *se
}

func (se *IdentifierVariable) display() string {
	return fmt.Sprintf("'%s'", se.value)
}

func CreateIdentifierVariable(value string) Variable {
	var result = IdentifierVariable{
		CommonVariable: CommonVariable{fType: TYPE_BOOL},
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

type Stack struct {
	// Storge of the stack, top element at index 0, bottom at length-1 (end of array)
	elts           []Variable
	onGoingSession bool
	listeners      []StackSessionListener
}

func CreateStack() *Stack {
	var s = Stack{}
	return &s
}

type StackSessionListener interface {
	SessionStart(s *Stack)
	SessionClose(s *Stack)
}

func CreateStackFromProto(reg *ActionRegistry, protoStack *protostack.Stack) (*Stack, error) {
	stack := CreateStack()
	for _, protoElt := range protoStack.Elements {
		variable, err := CreateVariableFromProto(reg, protoElt)
		if err != nil {
			return nil, err
		}
		stack.elts = append(stack.elts, variable)
	}
	return stack, nil
}

func CreateProtoFromStack(stack *Stack) (*protostack.Stack, error) {
	protoStack := &protostack.Stack{}
	for _, variable := range stack.elts {
		protoVar, err := CreateProtoFromVariable(variable)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal variable %w", err)
		}
		protoStack.Elements = append(protoStack.Elements, protoVar)
	}
	return protoStack, nil
}

func (s *Stack) Size() int {
	return len(s.elts)
}

/*
func (s *Stack) typeAt(l int) (Type, error) {
	if l < s.Size() {
		return (s.elts[len(s.elts)-l-1]).getType(), nil
	}
	return -1, fmt.Errorf("no elt at %d", l)
}
*/

func (s *Stack) IsEmpty() bool {
	return len(s.elts) == 0
}

func (s *Stack) Pop() (Variable, error) {
	if s.IsEmpty() {
		return nil, fmt.Errorf("empty stack")
	} else {
		index := len(s.elts) - 1
		result := s.elts[index]
		s.elts = s.elts[:index]
		return result, nil
	}
}

func (s *Stack) PopN(n int) ([]Variable, error) {
	if n == 0 {
		return []Variable{}, nil
	} else if s.Size() < n {
		return nil, fmt.Errorf("stack contains %d elements but %d were needed", s.Size(), n)
	} else {
		index := len(s.elts)
		result := make([]Variable, n)
		copy(result, s.elts[index-n:index])
		s.elts = s.elts[0 : index-n]
		return result, nil
	}
}

func (s *Stack) PeekN(n int) ([]Variable, error) {
	if n == 0 {
		return []Variable{}, nil
	} else if s.Size() < n {
		return nil, fmt.Errorf("stack contains %d elements but %d were needed", s.Size(), n)
	} else {
		index := len(s.elts)
		result := make([]Variable, n)
		// this copy is a bit conservative (operations could modify the slice we give them)
		copy(result, s.elts[index-n:index])
		return result, nil
	}
}

func (s *Stack) Get(level int) (Variable, error) {
	if level < s.Size() {
		return s.elts[len(s.elts)-level-1], nil
	} else {
		return nil, fmt.Errorf("Level %d does exist in stack of size %d", level, s.Size())
	}
}

func (s *Stack) Push(elt Variable) {
	s.elts = append(s.elts, elt)
	// fmt.Printf("After Push : len = %d\n", len(s.elts))
}

func (s *Stack) PushN(elts []Variable) {
	s.elts = append(s.elts, elts...)
}

func (s *Stack) StartSession() error {

	if s.onGoingSession {
		return fmt.Errorf("session already ongoing")
	}
	s.onGoingSession = true
	for _, listener := range s.listeners {
		listener.SessionStart(s)
	}
	return nil
}

func (s *Stack) CloseSession() error {
	if !s.onGoingSession {
		return fmt.Errorf("no ongoing session")
	}
	for _, listener := range s.listeners {
		listener.SessionClose(s)
	}
	s.onGoingSession = false
	return nil
}

type StackSavingListener struct {
	stackDataFolder string
}

func (sl *StackSavingListener) SessionStart(s *Stack) {

}

func (sl *StackSavingListener) SessionClose(s *Stack) {
	protoStack, err := CreateProtoFromStack(s)
	if err != nil {
		//TODO log error
		return
	}

	protoStackBytes, err := proto.Marshal(protoStack)
	if err != nil {
		//TODO log error
		return
	}
	err = os.WriteFile(sl.stackDataFolder, protoStackBytes, 0644)
	if err != nil {
		//TODO log error
		return
	}
}

func CreateSaveOnDiskStack(stackSavingPath string) *Stack {
	var stack *Stack
	file, err := os.ReadFile(stackSavingPath)
	if err != nil {
		stack = CreateStack()
	} else {
		protoStack := &protostack.Stack{}
		err = proto.Unmarshal(file, protoStack)
		if err != nil {
			stack = CreateStack()
		} else {
			stack, err = CreateStackFromProto(Registry, protoStack)
			if err != nil {
				stack = CreateStack()
			}
		}
	}
	saveStackSessionListener := &StackSavingListener{stackDataFolder: stackSavingPath}
	stack.listeners = append(stack.listeners, saveStackSessionListener)
	return stack
}
