package rcalc

import (
	"github.com/shopspring/decimal"
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
	if err != nil {
		t.Error(err)
	} else {
		i3, err := stack.Pop()
		if err != nil {
			t.Error(err)
		} else {
			if !i3.asNumericElt().value.Equals(decimal.NewFromInt(8)) {
				t.Errorf("3+5 should be 8 and not %d", i3.asNumericElt().value.IntPart())
			}
		}
	}

}
