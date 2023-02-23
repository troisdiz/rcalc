package rcalc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDisplayVariables(t *testing.T) {

	tests := []struct {
		variable        Variable
		expectedDisplay string
	}{
		{
			variable:        CreateProgramVariable([]Action{}),
			expectedDisplay: "<<  >>",
		},
		{
			variable: CreateProgramVariable([]Action{
				//&VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(4)},
				&VariableDeclarationActionDesc{
					varNames:           []string{"a"},
					variableToEvaluate: CreateProgramVariable([]Action{}),
				}}),
			expectedDisplay: "<< -> a <<  >> >>",
		}}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("Parse %02d", idx+1), func(t *testing.T) {
			programVariable := test.variable
			assert.Equal(t, test.expectedDisplay, programVariable.display())
		})
	}
}
