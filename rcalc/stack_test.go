package rcalc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyStack(t *testing.T) {
	var s Stack = CreateStack()
	assert.True(t, s.IsEmpty(), "Stack should be empty on creation")
}

func TestPeekN(t *testing.T) {
	var s Stack = CreateStack()
	se1 := CreateNumericVariableFromInt(2)
	s.Push(se1)
	se2 := CreateNumericVariableFromInt(7)
	s.Push(se2)
	peeked, _ := s.PeekN(2)
	expected := []Variable{se1, se2}
	assert.False(t, s.IsEmpty(), "Stack should NOT be empty after Push and Pop (current size %d)", s.Size())

	if assert.Len(t, peeked, len(expected), "Peeked array length is %d and not %d", len(peeked), len(expected)) {
		for idx, se := range peeked {
			assert.Equal(t, expected[idx], se, "Peeked element at index %d is not the expected one (real = %v, expected, %v)", idx, se, expected[idx])
		}
	}
}

func TestPushAndPop(t *testing.T) {
	var s Stack = CreateStack()
	se := CreateNumericVariableFromInt(2)
	s.Push(se)
	popped, _ := s.Pop()

	assert.True(t, s.IsEmpty(), "Stack should be empty after Push and Pop")
	assert.Equal(t, Variable(se), popped, "Popped elt is not the inserted one")
}

func TestPushAndPopN(t *testing.T) {
	var s Stack = CreateStack()
	var se1, se2 Variable
	se1 = CreateNumericVariable(decimal.NewFromInt(2))
	s.Push(se1)
	se2 = CreateBooleanVariable(true)
	s.Push(se2)
	popped, _ := s.PopN(2)

	assert.True(t, s.IsEmpty(), "Stack should be empty after Push and Pop")
	assert.Equal(t, 2, len(popped), "Popped has not the right length (%d instead of 2)", len(popped))
}

func TestPushAndPopNAndSize(t *testing.T) {
	var s Stack = CreateStack()
	var se1, se2 Variable
	se1 = CreateNumericVariable(decimal.NewFromInt(2))
	s.Push(se1)
	se2 = CreateBooleanVariable(true)
	s.Push(se2)
	popped, _ := s.PopN(1)

	assert.False(t, s.IsEmpty(), "Stack should not be empty after Push 2 and Pop 1")
	assert.Equal(t, 1, len(popped), "Popped has not the right length (%d)", len(popped))
}

func TestPushAndSize(t *testing.T) {
	var s Stack = CreateStack()
	se := CreateNumericVariableFromInt(2)
	s.Push(se)
	// fmt.Printf("Size after 1 Push %d / %d\n", s.Size(), len(s.elts))
	se2 := CreateNumericVariableFromInt(2)
	s.Push(se2)

	assert.Equal(t, 2, s.Size(), "Stack Size must be 2 and id %d", s.Size())
}

func TestDisplayStack(t *testing.T) {
	var s Stack = CreateStack()
	se := CreateNumericVariable(decimal.NewFromInt(2))
	s.Push(se)
	fmt.Printf("Size after 1 Push %d / %d\n", s.Size(), len(s.elts))
	se2 := CreateNumericVariable(decimal.NewFromInt(3))
	s.Push(se2)
	DisplayStack(s, "", 4, false)
}

func TestNumericStackEltType(t *testing.T) {
	bse := CreateNumericVariable(decimal.NewFromInt(5))
	assert.Equal(t, TYPE_NUMERIC, bse.getType(), "Type should be %d and is %d", TYPE_NUMERIC, bse.getType())
}

func TestBooleanStackEltType(t *testing.T) {
	bse := CreateBooleanVariable(true)
	assert.Equal(t, TYPE_BOOL, bse.getType(), "Type should be %d and is %d", TYPE_BOOL, bse.getType())
}
