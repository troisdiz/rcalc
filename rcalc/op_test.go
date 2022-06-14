package rcalc

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddApply(t *testing.T) {
	var i1 = CreateNumericVariableFromInt(3)

	var i2 = CreateNumericVariableFromInt(5)

	stack := CreateStack()
	stack.Push(i1)
	stack.Push(i2)
	err := addOp.Apply(nil, &stack)
	if assert.NoError(t, err) {
		i3, err := stack.Pop()
		if assert.NoError(t, err) {
			assert.Equal(t, decimal.NewFromInt(8), i3.asNumericVar().value, "3+5 should be 8 and not %d", i3.asNumericVar().value.IntPart())
		}
	}
}
