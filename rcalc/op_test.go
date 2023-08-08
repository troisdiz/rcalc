package rcalc

import (
	"fmt"
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
	runtimeContext := CreateRuntimeContext(nil, stack)
	err := addOp.Apply(runtimeContext)
	if assert.NoError(t, err) {
		i3, err := stack.Pop()
		if assert.NoError(t, err) {
			assert.Equal(t, decimal.NewFromInt(8), i3.asNumericVar().value, "3+5 should be 8 and not %d", i3.asNumericVar().value.IntPart())
		}
	}
}

func TestCheckGenOk(t *testing.T) {
	i1 := CreateNumericVariableFromInt(5)
	i2 := CreateBooleanVariable(true)

	ok, err := CheckGen([]Type{TYPE_NUMERIC, TYPE_BOOL})(i1, i2)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestCheckGenGeneric(t *testing.T) {
	i1 := CreateNumericVariableFromInt(5)
	i2 := CreateBooleanVariable(true)

	ok, err := CheckGen([]Type{TYPE_NUMERIC, TYPE_GENERIC})(i1, i2)
	assert.True(t, ok)
	assert.NoError(t, err)
}

type instrumentedExpandableFunctionFactory struct {
	callingCount int
}

func (i *instrumentedExpandableFunctionFactory) GetApplyFunc() OperationApplyFn {
	return func(system System, elts ...Variable) []Variable {
		i.callingCount += 1
		fmt.Printf("Applyfunc called\n")
		return []Variable{CreateListVariable([]Variable{elts[0], elts[1]})}
	}
}

func (i *instrumentedExpandableFunctionFactory) GetApplyFunc22() OperationApplyFn {
	return func(system System, elts ...Variable) []Variable {
		i.callingCount += 1
		fmt.Printf("Applyfunc called\n")
		return []Variable{elts[0], elts[1]}
	}
}

func TestExpandableApply(t *testing.T) {

	t.Run("NotExpandedCall", func(t *testing.T) {

		funcBuilder := new(instrumentedExpandableFunctionFactory)

		expandableTestOp := NewExpandableOperationDesc("testExpandableOp", 2, CheckNoop, 1, funcBuilder.GetApplyFunc())

		var i1 = CreateNumericVariableFromInt(3)
		var i2 = CreateNumericVariableFromInt(5)

		stack := CreateStack()
		stack.Push(i1)
		stack.Push(i2)
		runtimeContext := CreateRuntimeContext(nil, stack)
		err := expandableTestOp.Apply(runtimeContext)
		if assert.NoError(t, err) {
			i3, err := stack.Pop()
			if assert.NoError(t, err) {
				//assert.Equal(t, decimal.NewFromInt(8), i3.asNumericVar().value, "3+5 should be 8 and not %d", i3.asNumericVar().value.IntPart())
				assert.NotNilf(t, i3, "")
				assert.Equal(t, CreateListVariable([]Variable{CreateNumericVariableFromInt(3), CreateNumericVariableFromInt(5)}), i3)
			}
		}
	})
	t.Run("ExpandedCall 2-1", func(t *testing.T) {

		funcBuilder := new(instrumentedExpandableFunctionFactory)

		expandableTestOp := NewExpandableOperationDesc("testExpandableOp", 2, CheckNoop, 1, funcBuilder.GetApplyFunc())

		var i1 = CreateListVariable([]Variable{CreateNumericVariableFromInt(3), CreateNumericVariableFromInt(13)})
		var i2 = CreateListVariable([]Variable{CreateNumericVariableFromInt(5), CreateNumericVariableFromInt(15)})

		stack := CreateStack()
		stack.Push(i1)
		stack.Push(i2)
		runtimeContext := CreateRuntimeContext(nil, stack)
		err := expandableTestOp.Apply(runtimeContext)
		if assert.NoError(t, err) {
			i3, err := stack.Pop()
			fmt.Println()
			fmt.Printf("%s\n%s\n-> %s", i1.display(), i2.display(), i3.display())
			if assert.NoError(t, err) {
				assert.Equal(t, TYPE_LIST, i3.getType())
				i3List := i3.asListVar()
				assert.Equal(t, 2, i3List.Size())
				assert.Equal(t, TYPE_LIST, i3List.items[0].getType())
				assert.Equal(t, TYPE_LIST, i3List.items[1].getType())

				i3List0 := i3List.items[0].asListVar()
				i3List1 := i3List.items[1].asListVar()
				assert.Equal(t, 2, i3List0.Size(), i3List0.display())
				assert.Equal(t, 2, i3List1.Size(), i3List1.display())

				i3List00 := i3List0.items[0].asNumericVar()
				i3List01 := i3List0.items[1].asNumericVar()
				i3List10 := i3List1.items[0].asNumericVar()
				i3List11 := i3List1.items[1].asNumericVar()

				assert.True(t, decimal.NewFromInt(3).Equal(i3List00.value), i3List00.display())
				assert.True(t, decimal.NewFromInt(5).Equal(i3List01.value), i3List01.display())
				assert.True(t, decimal.NewFromInt(13).Equal(i3List10.value), i3List10.display())
				assert.True(t, decimal.NewFromInt(15).Equal(i3List11.value), i3List11.display())
			}
		}
	})
	t.Run("ExpandedCall 2-2", func(t *testing.T) {

		funcBuilder := new(instrumentedExpandableFunctionFactory)

		expandableTestOp := NewExpandableOperationDesc("testExpandableOp", 2, CheckNoop, 2, funcBuilder.GetApplyFunc22())

		var i1 = CreateListVariable([]Variable{CreateNumericVariableFromInt(3), CreateNumericVariableFromInt(13)})
		var i2 = CreateListVariable([]Variable{CreateNumericVariableFromInt(5), CreateNumericVariableFromInt(15)})

		stack := CreateStack()
		stack.Push(i1)
		stack.Push(i2)
		runtimeContext := CreateRuntimeContext(nil, stack)
		err := expandableTestOp.Apply(runtimeContext)
		if assert.NoError(t, err) {

			if assert.Equal(t, 2, stack.Size()) {

				i3, err := stack.Pop()
				fmt.Println()
				fmt.Printf("%s\n%s\n-> %s\n", i1.display(), i2.display(), i3.display())
				if assert.NoError(t, err) {
					assert.Equal(t, TYPE_LIST, i3.getType())
					i3List := i3.asListVar()
					assert.Equal(t, 2, i3List.Size())
					assert.Equal(t, TYPE_LIST, i3List.items[0].getType(), "Expected TYPE_LIST(6) and got other type")
					assert.Equal(t, TYPE_LIST, i3List.items[1].getType(), "Expected TYPE_LIST(6) and got other type")

					i3List0 := i3List.items[0].asListVar()
					i3List1 := i3List.items[1].asListVar()
					assert.Equal(t, 2, i3List0.Size(), i3List0.display())
					assert.Equal(t, 2, i3List1.Size(), i3List1.display())

					i3List00 := i3List0.items[0].asNumericVar()
					i3List01 := i3List0.items[1].asNumericVar()
					i3List10 := i3List1.items[0].asNumericVar()
					i3List11 := i3List1.items[1].asNumericVar()

					assert.True(t, decimal.NewFromInt(3).Equal(i3List00.value), i3List00.display())
					assert.True(t, decimal.NewFromInt(5).Equal(i3List01.value), i3List01.display())
					assert.True(t, decimal.NewFromInt(13).Equal(i3List10.value), i3List10.display())
					assert.True(t, decimal.NewFromInt(15).Equal(i3List11.value), i3List11.display())
				}
			}
		}
	})
}
