package rcalc

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddApply(t *testing.T) {
	var i1 = NumericStackElt{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(3),
	}

	var i2 = NumericStackElt{
		fType: TYPE_NUMERIC,
		value: decimal.NewFromInt(5),
	}

	stack := CreateStack()
	stack.Push(&i1)
	stack.Push(&i2)
	err := addOp.Apply(nil, &stack)
	if assert.NoError(t, err) {
		i3, err := stack.Pop()
		if assert.NoError(t, err) {
			assert.Equal(t, decimal.NewFromInt(8), i3.asNumericElt().value, "3+5 should be 8 and not %d", i3.asNumericElt().value.IntPart())
		}
	}
}
