package rcalc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
