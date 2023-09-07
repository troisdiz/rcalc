package rcalc

import (
	"fmt"
	"os"
	"slices"

	"google.golang.org/protobuf/proto"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

type Type int

const (
	TYPE_GENERIC  Type = 0
	TYPE_NUMERIC  Type = 1
	TYPE_BOOL     Type = 2
	TYPE_STR      Type = 3
	TYPE_ALG_EXPR Type = 4
	TYPE_PROGRAM  Type = 5
	TYPE_LIST     Type = 6
	// TYPE_VECTOR     Type = 7
)

type Variable interface {
	getType() Type
	asNumericVar() *NumericVariable
	asBooleanVar() *BooleanVariable
	asIdentifierVar() *AlgebraicExpressionVariable
	asProgramVar() *ProgramVariable
	asListVar() *ListVariable
	display() string
	String() string
}

type CommonVariable struct {
	fType Type
}

func (se *CommonVariable) getType() Type {
	return se.fType
}

func (se *CommonVariable) asNumericVar() *NumericVariable {
	panic("This is not a Numeric variable")
}

func (se *CommonVariable) asBooleanVar() *BooleanVariable {
	panic("This is not a Boolean variable")
}

func (se *CommonVariable) asIdentifierVar() *AlgebraicExpressionVariable {
	panic("This is not an Identifier variable")
}

func (se *CommonVariable) asProgramVar() *ProgramVariable {
	panic("This is not a Program variable")
}

func (se *CommonVariable) asListVar() *ListVariable {
	panic("This is not a List variable")
}

func (se *CommonVariable) String() string {
	return fmt.Sprintf("[CommonVariable] t=%d", se.fType)
}

type Stack struct {
	// Store of the stack, top element at index 0, bottom at length-1 (end of array)
	current []Variable
	// Last item = more recent
	history        [][]Variable
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

func CreateStackFromProto(reg *ActionRegistry, protoStackWithHistory *protostack.StackWithHistory) (*Stack, error) {
	stack := CreateStack()
	for _, protoElt := range protoStackWithHistory.Current.Elements {
		variable, err := CreateVariableFromProto(reg, protoElt)
		if err != nil {
			return nil, err
		}
		stack.current = append(stack.current, variable)
	}
	for _, protoHistoricStack := range protoStackWithHistory.History {
		var historicStack []Variable
		for _, protoElt := range protoHistoricStack.Elements {
			variable, err := CreateVariableFromProto(reg, protoElt)
			if err != nil {
				return nil, err
			}
			historicStack = append(historicStack, variable)
		}
		stack.history = append(stack.history, historicStack)
	}
	return stack, nil
}

func CreateProtoFromStack(stack *Stack) (*protostack.StackWithHistory, error) {
	protoStackWithHistory := &protostack.StackWithHistory{}
	protoStackWithHistory.Current = &protostack.Stack{}
	for _, variable := range stack.current {
		protoVar, err := CreateProtoFromVariable(variable)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal variable %w", err)
		}
		protoStackWithHistory.Current.Elements = append(protoStackWithHistory.Current.Elements, protoVar)
	}
	for _, historicStack := range stack.history {
		protoStack := &protostack.Stack{}
		for _, variable := range historicStack {
			protoVar, err := CreateProtoFromVariable(variable)
			if err != nil {
				return nil, fmt.Errorf("cannot marshal variable %w", err)
			}
			protoStack.Elements = append(protoStack.Elements, protoVar)
		}
		protoStackWithHistory.History = append(protoStackWithHistory.History, protoStack)
	}
	return protoStackWithHistory, nil
}

func (s *Stack) Size() int {
	return len(s.current)
}

func (s *Stack) HistoryLength() int {
	return len(s.history)
}

func (s *Stack) IsEmpty() bool {
	return len(s.current) == 0
}

func (s *Stack) Pop() (Variable, error) {
	if s.IsEmpty() {
		return nil, fmt.Errorf("empty stack")
	} else {
		index := s.Size() - 1
		result := s.current[index]
		s.current = (s.current)[:index]
		return result, nil
	}
}

func (s *Stack) PopN(n int) ([]Variable, error) {
	if n == 0 {
		return []Variable{}, nil
	} else if s.Size() < n {
		return nil, fmt.Errorf("stack contains %d elements but %d were needed", s.Size(), n)
	} else {
		index := s.Size()
		result := make([]Variable, n)
		copy(result, (s.current)[index-n:index])
		s.current = s.current[0 : index-n]
		return result, nil
	}
}

func (s *Stack) PeekN(n int) ([]Variable, error) {
	if n == 0 {
		return []Variable{}, nil
	} else if s.Size() < n {
		return nil, fmt.Errorf("stack contains %d elements but %d were needed", s.Size(), n)
	} else {
		index := len(s.current)
		result := make([]Variable, n)
		// this copy is a bit conservative (operations could modify the slice we give them)
		copy(result, (s.current)[index-n:index])
		return result, nil
	}
}

func (s *Stack) Get(level int) (Variable, error) {
	if level < s.Size() {
		return s.current[len(s.current)-level-1], nil
	} else {
		return nil, fmt.Errorf("level %d does not exist in stack of size %d", level, s.Size())
	}
}

func (s *Stack) Push(elt Variable) {
	s.current = append(s.current, elt)
	// fmt.Printf("After Push : len = %d\n", len(s.history))
}

func (s *Stack) PushN(elts []Variable) {
	s.current = append(s.current, elts...)
}

func (s *Stack) StartSession() error {

	if s.onGoingSession {
		return fmt.Errorf("session already ongoing")
	}
	s.onGoingSession = true
	for _, listener := range s.listeners {
		listener.SessionStart(s)
	}
	var newStack []Variable
	copy(newStack, s.current)
	s.history = append(s.history, newStack)
	return nil
}

func (s *Stack) CloseSession() error {
	if !s.onGoingSession {
		return fmt.Errorf("no ongoing session")
	}
	for _, listener := range slices.Backward(s.listeners) {
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
		GetLogger().Debugf("Error saving stack: %v", err)
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
		protoStackHistory := &protostack.StackWithHistory{}
		err = proto.Unmarshal(file, protoStackHistory)
		if err != nil {
			stack = CreateStack()
		} else {
			stack, err = CreateStackFromProto(Registry, protoStackHistory)
			if err != nil {
				stack = CreateStack()
			}
		}
	}
	saveStackSessionListener := &StackSavingListener{stackDataFolder: stackSavingPath}
	stack.listeners = append(stack.listeners, saveStackSessionListener)
	return stack
}
