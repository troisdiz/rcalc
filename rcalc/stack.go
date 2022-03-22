package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type Type int

const (
	TYPE_NUMERIC Type = 0
	TYPE_STR     Type = 1
)

type NumericStackElt struct {
	fType Type
	value decimal.Decimal
}

func (se *NumericStackElt) String() string {
	return fmt.Sprintf("NumericStackElt(%v)", se.value)
}

func CreateNumericStackElt(value decimal.Decimal) StackElt {
	var result = NumericStackElt{
		fType: TYPE_NUMERIC,
		value: value,
	}
	return &result
}

func (se *NumericStackElt) asNumericElt() NumericStackElt {
	return *se
}

func (se *NumericStackElt) getType() Type {
	return 0
}

func (se *NumericStackElt) display() string {
	return se.value.String()
}

type StackElt interface {
	getType() Type
	asNumericElt() NumericStackElt
	display() string
}

type Stack struct {
	elts []StackElt
}

func CreateStack() Stack {
	var s = Stack{}
	return s
}

func (s *Stack) Size() int {
	return len(s.elts)
}

func (s *Stack) typeAt(l int) (Type, error) {
	if l < s.Size() {
		return (s.elts[len(s.elts)-l-1]).getType(), nil
	}
	return -1, fmt.Errorf("no elt at %d", l)
}

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
