package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type Type int

const (
	TYPE_NUMERIC Type = 0
	TYPE_BOOL    Type = 1
	TYPE_STR     Type = 2
)

type Variable interface {
	getType() Type
	asNumericVar() NumericVariable
	asBooleanVar() BooleanVariable
	display() string
}

type NumericVariable struct {
	fType Type
	value decimal.Decimal
}

func (se *NumericVariable) String() string {
	return fmt.Sprintf("NumericVariable(%v)", se.value)
}

func (se *NumericVariable) asNumericVar() NumericVariable {
	return *se
}

func (se *NumericVariable) asBooleanVar() BooleanVariable {
	panic("This is a Numeric and not boolean element")
}

func (se *NumericVariable) getType() Type {
	return se.fType
}

func (se *NumericVariable) display() string {
	return se.value.String()
}

func CreateNumericVariable(value decimal.Decimal) Variable {
	var result = NumericVariable{
		fType: TYPE_NUMERIC,
		value: value,
	}
	return &result
}

func CreateNumericVariableFromInt(value int) Variable {
	var result = NumericVariable{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(int64(value)),
	}
	return &result
}

type BooleanVariable struct {
	fType Type
	value bool
}

func (se *BooleanVariable) String() string {
	return fmt.Sprintf("BooleanVariable(%v) type = %d", se.value, se.fType)
}

func (se *BooleanVariable) asNumericVar() NumericVariable {
	panic("This is a Boolean and not Numeric element")
}

func (se *BooleanVariable) asBooleanVar() BooleanVariable {
	return *se
}

func (se *BooleanVariable) getType() Type {
	return se.fType
}

func (se *BooleanVariable) display() string {
	return fmt.Sprintf("%t", se.value)
}

func CreateBooleanVariable(value bool) Variable {
	var result = BooleanVariable{
		fType: TYPE_BOOL,
		value: value,
	}
	return &result
}

type Stack struct {
	// Storge of the stack, top element at index 0, bottom at length-1 (end of array)
	elts []Variable
}

func CreateStack() Stack {
	var s = Stack{}
	return s
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
