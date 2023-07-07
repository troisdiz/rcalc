package rcalc

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpsPow(t *testing.T) {

	var x = CreateNumericVariable(decimal.NewFromInt(2))
	var y = CreateNumericVariable(decimal.NewFromInt(3))

	stack := CreateStack()
	stack.Push(x)
	stack.Push(y)

	system := CreateSystemInstance()
	runtimeContext := CreateRuntimeContext(system, stack)

	var err = powOp.Apply(runtimeContext)
	if assert.NoError(t, err) {
		var result Variable
		result, err = stack.Pop()
		if assert.NoError(t, err) {
			assert.Equal(t, decimal.NewFromInt(8), result.asNumericVar().value)
		}
	}
}

func TestOpsSub(t *testing.T) {

	var x = CreateNumericVariable(decimal.NewFromInt(2))
	var y = CreateNumericVariable(decimal.NewFromInt(3))

	stack := CreateStack()
	stack.Push(x)
	stack.Push(y)

	system := CreateSystemInstance()
	runtimeContext := CreateRuntimeContext(system, stack)

	var err = subOp.Apply(runtimeContext)
	if assert.NoError(t, err) {
		var result Variable
		result, err = stack.Pop()
		if assert.NoError(t, err) {
			assert.Equal(t, decimal.NewFromInt(-1), result.asNumericVar().value)
		}
	}
}

func TestOpsDiv(t *testing.T) {

	var x = CreateNumericVariable(decimal.NewFromInt(27))
	var y = CreateNumericVariable(decimal.NewFromInt(3))

	stack := CreateStack()
	stack.Push(x)
	stack.Push(y)

	system := CreateSystemInstance()
	runtimeContext := CreateRuntimeContext(system, stack)

	var err = divOp.Apply(runtimeContext)
	if assert.NoError(t, err) {
		var result Variable
		result, err = stack.Pop()
		if assert.NoError(t, err) {
			assert.True(t, decimal.NewFromInt(9).Equal(result.asNumericVar().value))
		}
	}
}

func TestCrDirAction(t *testing.T) {

	myFolderName := "MyFolder"
	var id1 = CreateAlgebraicExpressionVariable(myFolderName, nil)

	stack := CreateStack()
	stack.Push(id1)

	system := CreateSystemInstance()
	runtimeContext := CreateRuntimeContext(system, stack)
	err := crdirAct.Apply(runtimeContext)
	assert.NoError(t, err, "Creation of folder should work")
	rootFolder := system.Memory().getRoot()
	subFolders := rootFolder.subFolders
	found := false
	for _, f := range subFolders {
		if f.name == myFolderName {
			found = true
			break
		}
	}
	assert.True(t, found)

}
