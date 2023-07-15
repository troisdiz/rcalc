package rcalc

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ArithmeticTestSuite struct {
	suite.Suite
}

func (suite *ArithmeticTestSuite) testOperation(operation OperationDesc, inputVars []Variable, expectedOutputVars []*NumericVariable) {
	stack := CreateStack()
	for _, inputVar := range inputVars {
		stack.Push(inputVar)

	}
	system := CreateSystemInstance()
	runtimeContext := CreateRuntimeContext(system, stack)

	var err = operation.Apply(runtimeContext)
	if assert.NoError(suite.T(), err) {
		var results []Variable

		results, err = stack.PopN(len(expectedOutputVars))
		if assert.NoError(suite.T(), err) {
			for i, expectedOutputVar := range expectedOutputVars {
				assert.True(suite.T(), expectedOutputVar.Equals(results[i].asNumericVar()),
					"Expected %s, got %s",
					expectedOutputVar.value.String(),
					results[i].asNumericVar().value.String())
			}
		}
	}
}

func TestArithmeticTestSuite(t *testing.T) {

	suite.Run(t, new(ArithmeticTestSuite))
}

func (suite *ArithmeticTestSuite) TestOpsPow() {

	suite.testOperation(powOp,
		[]Variable{
			CreateNumericVariable(decimal.NewFromInt(2)),
			CreateNumericVariable(decimal.NewFromInt(3)),
		},
		[]*NumericVariable{
			CreateNumericVariable(decimal.NewFromInt(8)).asNumericVar(),
		},
	)
}

func (suite *ArithmeticTestSuite) TestOpsSub() {
	suite.testOperation(subOp,
		[]Variable{
			CreateNumericVariable(decimal.NewFromInt(2)),
			CreateNumericVariable(decimal.NewFromInt(3)),
		},
		[]*NumericVariable{
			CreateNumericVariable(decimal.NewFromInt(-1)).asNumericVar(),
		},
	)
}

func (suite *ArithmeticTestSuite) TestOpsDiv() {
	suite.testOperation(divOp,
		[]Variable{
			CreateNumericVariable(decimal.NewFromInt(27)),
			CreateNumericVariable(decimal.NewFromInt(3)),
		},
		[]*NumericVariable{
			CreateNumericVariable(decimal.NewFromInt(9)).asNumericVar(),
		},
	)
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
