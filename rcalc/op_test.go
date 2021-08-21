package rcalc

import "testing"

func TestAddApply(t *testing.T) {
	var i1 = NumericStackElt{
		fType: TYPE_NUMERIC,
		value: 3,
	}

	var i2 = NumericStackElt{
		fType: TYPE_NUMERIC,
		value: 5,
	}

	i3 := ADD_OP.Apply(nil, &i1, &i2).asNumericElt()
	if i3.value != 8 {
		t.Errorf("3+5 should be 8 and not %d", i3.value)
	}
}
