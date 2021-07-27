package rcalc

import (
	"fmt"
	"testing"
)

func TestEmptyStack(t *testing.T) {
	var s Stack = CreateStack()
	if !s.IsEmpty() {
		t.Error("Stack should be empty on creation")
	}
}

func TestPushAndPop(t *testing.T) {
	var s Stack = CreateStack()
	se := IntStackElt{
		fType: TYPE_INT,
		value: 2,
	}
	s.Push(&se)
	popped, _ := s.Pop()

	if !s.IsEmpty() {
		t.Error("Stack should be empty after Push and Pop")
	}
	if popped != StackElt(&se) {
		t.Error("Popped elt is not the inserted one")
	}
}

func TestPushAndSize(t *testing.T) {
	var s Stack = CreateStack()
	se := IntStackElt{
		fType: TYPE_INT,
		value: 2,
	}
	s.Push(&se)
	fmt.Printf("Size after 1 Push %d / %d\n", s.Size(), len(s.elts))
	se2 := IntStackElt{
		fType: TYPE_INT,
		value: 2,
	}
	s.Push(&se2)

	if s.Size() != 2 {
		t.Errorf("Stack Size must be 2 and id %d", s.Size())
	}
}

func TestDisplayStack(t *testing.T) {
	var s Stack = CreateStack()
	se := CreateIntStackElt(2)
	s.Push(se)
	fmt.Printf("Size after 1 Push %d / %d\n", s.Size(), len(s.elts))
	se2 := CreateIntStackElt(3)
	s.Push(se2)
	DisplayStack(s, 4)
}
