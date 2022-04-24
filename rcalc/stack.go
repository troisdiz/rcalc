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

type StackElt interface {
	getType() Type
	asNumericElt() NumericStackElt
	asBooleanElt() BooleanStackElt
	display() string
}

type NumericStackElt struct {
	fType Type
	value decimal.Decimal
}

func (se *NumericStackElt) String() string {
	return fmt.Sprintf("NumericStackElt(%v)", se.value)
}

func (se *NumericStackElt) asNumericElt() NumericStackElt {
	return *se
}

func (se *NumericStackElt) asBooleanElt() BooleanStackElt {
	panic("This is a Numeric and not boolean element")
}

func (se *NumericStackElt) getType() Type {
	return se.fType
}

func (se *NumericStackElt) display() string {
	return se.value.String()
}

func CreateNumericStackElt(value decimal.Decimal) StackElt {
	var result = NumericStackElt{
		fType: TYPE_NUMERIC,
		value: value,
	}
	return &result
}

type BooleanStackElt struct {
	fType Type
	value bool
}

func (se *BooleanStackElt) String() string {
	return fmt.Sprintf("BooleanStackElt(%v) type = %d", se.value, se.fType)
}

func (se *BooleanStackElt) asNumericElt() NumericStackElt {
	panic("This is a Boolean and not Numeric element")
}

func (se *BooleanStackElt) asBooleanElt() BooleanStackElt {
	return *se
}

func (se *BooleanStackElt) getType() Type {
	return se.fType
}

func (se *BooleanStackElt) display() string {
	return fmt.Sprintf("%t", se.value)
}

func CreateBooleanStackElt(value bool) StackElt {
	var result = BooleanStackElt{
		fType: TYPE_BOOL,
		value: value,
	}
	return &result
}

type Stack struct {
	// Storge of the stack, top element at index 0, bottom at length-1 (end of array)
	elts []StackElt
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

func (s *Stack) Pop() (StackElt, error) {
	if s.IsEmpty() {
		return nil, fmt.Errorf("empty stack")
	} else {
		index := len(s.elts) - 1
		result := s.elts[index]
		s.elts = s.elts[:index]
		return result, nil
	}
}

func (s *Stack) PopN(n int) ([]StackElt, error) {
	if n == 0 {
		return []StackElt{}, nil
	} else if s.Size() < n {
		return nil, fmt.Errorf("stack contains %d elements but %d were needed", s.Size(), n)
	} else {
		index := len(s.elts)
		result := make([]StackElt, n)
		copy(result, s.elts[index-n:index])
		s.elts = s.elts[0 : index-n]
		return result, nil
	}
}

func (s *Stack) PeekN(n int) ([]StackElt, error) {
	if n == 0 {
		return []StackElt{}, nil
	} else if s.Size() < n {
		return nil, fmt.Errorf("stack contains %d elements but %d were needed", s.Size(), n)
	} else {
		index := len(s.elts)
		result := make([]StackElt, n)
		// this copy is probably
		copy(result, s.elts[index-n:index])
		return result, nil
	}
}

func (s *Stack) Get(level int) (StackElt, error) {
	if level < s.Size() {
		return s.elts[len(s.elts)-level-1], nil
	} else {
		return nil, fmt.Errorf("Level %d does exist in stack of size %d", level, s.Size())
	}
}

func (s *Stack) Push(elt StackElt) {
	s.elts = append(s.elts, elt)
	// fmt.Printf("After Push : len = %d\n", len(s.elts))
}
