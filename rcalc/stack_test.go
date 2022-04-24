package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestEmptyStack(t *testing.T) {
	var s Stack = CreateStack()
	if !s.IsEmpty() {
		t.Error("Stack should be empty on creation")
	}
}

func TestPeekN(t *testing.T) {
	var s Stack = CreateStack()
	se1 := NumericStackElt{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(2),
	}
	s.Push(&se1)
	se2 := NumericStackElt{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(7),
	}
	s.Push(&se2)
	peeked, _ := s.PeekN(2)
	expected := []StackElt{&se1, &se2}
	if s.IsEmpty() {
		t.Errorf("Stack should NOT be empty after Push and Pop (current size %d)", s.Size())
	}

	if len(peeked) != len(expected) {
		t.Errorf("Peeked array length is %d and not %d", len(peeked), len(expected))
	} else {
		for idx, se := range peeked {
			if se != expected[idx] {
				t.Errorf("Peeked element at index %d is not the expected one (real = %v, expected, %v)", idx, se, expected[idx])
			}
		}
	}
}

func TestPushAndPop(t *testing.T) {
	var s Stack = CreateStack()
	se := NumericStackElt{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(2),
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

func TestPushAndPopN(t *testing.T) {
	var s Stack = CreateStack()
	var se1, se2 StackElt
	se1 = CreateNumericStackElt(decimal.NewFromInt(2))
	s.Push(se1)
	se2 = CreateBooleanStackElt(true)
	s.Push(se2)
	popped, _ := s.PopN(2)

	if !s.IsEmpty() {
		t.Error("Stack should be empty after Push and Pop")
	}
	if len(popped) != 2 {
		t.Errorf("Popped has not the right length (%d instead of 2)", len(popped))
	}

}

func TestPushAndPopNAndSize(t *testing.T) {
	var s Stack = CreateStack()
	var se1, se2 StackElt
	se1 = CreateNumericStackElt(decimal.NewFromInt(2))
	s.Push(se1)
	se2 = CreateBooleanStackElt(true)
	s.Push(se2)
	popped, _ := s.PopN(1)

	if s.IsEmpty() {
		t.Error("Stack should not be empty after Push 2 and Pop 1")
	}
	if len(popped) != 1 {
		t.Error("Popped has not the right length")
	}
}

func TestPushAndSize(t *testing.T) {
	var s Stack = CreateStack()
	se := NumericStackElt{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(2),
	}
	s.Push(&se)
	fmt.Printf("Size after 1 Push %d / %d\n", s.Size(), len(s.elts))
	se2 := NumericStackElt{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(2),
	}
	s.Push(&se2)

	if s.Size() != 2 {
		t.Errorf("Stack Size must be 2 and id %d", s.Size())
	}
}

func TestDisplayStack(t *testing.T) {
	var s Stack = CreateStack()
	se := CreateNumericStackElt(decimal.NewFromInt(2))
	s.Push(se)
	fmt.Printf("Size after 1 Push %d / %d\n", s.Size(), len(s.elts))
	se2 := CreateNumericStackElt(decimal.NewFromInt(3))
	s.Push(se2)
	DisplayStack(s, "", 4)
}

func TestNumericStackEltType(t *testing.T) {
	bse := CreateNumericStackElt(decimal.NewFromInt(5))
	if bse.getType() != TYPE_NUMERIC {
		t.Errorf("Type should be %d and is %d", TYPE_NUMERIC, bse.getType())
	}
}

func TestBooleanStackEltType(t *testing.T) {
	bse := CreateBooleanStackElt(true)
	if bse.getType() != TYPE_BOOL {
		t.Errorf("Type should be %d and is %d", TYPE_BOOL, bse.getType())
	}
}
