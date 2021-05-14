package rcalc

import "testing"

func TestAddApply(t *testing.T) {
	var i1 = IntStackElt{
		fType: TYPE_INT,
		value: 3,
	}

	var i2 = IntStackElt{
		fType: TYPE_INT,
		value: 5,
	}

	i3 := ADD_OP.Apply(&i1, &i2).asIntElt()
	if i3.value != 8 {
		t.Errorf("3+5 should be 8 and not %d", i3.value)
	}
}
