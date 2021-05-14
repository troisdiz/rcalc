package rcalc

import (
	"fmt"
	"strconv"
)

type Type int

const (
	TYPE_INT Type = 0
	TYPE_STR Type = 1
)

type IntStackElt struct {
	fType Type
	value int
}

func CreateInStackElt(value int) *IntStackElt {
	var result = IntStackElt{
		fType: TYPE_INT,
		value: value,
	}
	return &result
}

func (se *IntStackElt) asIntElt() IntStackElt {
	return *se
}

func (se *IntStackElt) getType() Type {
	return 0
}

func (se *IntStackElt) display() string {
	return strconv.Itoa(se.value)
}

type StackElt interface {
	getType() Type
	asIntElt() IntStackElt
	display()  string
}


type Stack struct {
	elts []StackElt
}

func Create() Stack {
	var s = Stack{}
	return s
}

func (s *Stack) Size() int {
	return len(s.elts)
}

func (s *Stack) typeAt(l int) (Type, error)  {
	if l < s.Size() {
		return (s.elts[len(s.elts) - l -1]).getType(), nil
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

func (s *Stack) Get(level int) (StackElt, error)  {
	if level < s.Size() {
		return s.elts[len(s.elts)-level-1], nil
	} else {
		return nil, fmt.Errorf("Level %d does exist in stack of size %d", level, s.Size())
	}
}

func (s *Stack) Push(elt StackElt)  {
	s.elts = append(s.elts, elt)
	fmt.Printf("After Push : len = %d\n", len(s.elts))
}
