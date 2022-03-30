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
