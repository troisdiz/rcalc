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

	i3 := addOp.Apply(nil, &i1, &i2)[0].asNumericElt()
	if !i3.value.Equals(decimal.NewFromInt(8)) {
		t.Errorf("3+5 should be 8 and not %d", i3.value.IntPart())
	}
}
