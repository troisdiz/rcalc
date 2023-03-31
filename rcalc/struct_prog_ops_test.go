package rcalc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartNextLoop(t *testing.T) {
	var start = CreateNumericVariableFromInt(1)
	var end = CreateNumericVariableFromInt(3)
	var one = CreateNumericVariableFromInt(1)

	stack := CreateStack()
	stack.Push(one)
	stack.Push(start)
	stack.Push(end)

	system := CreateSystemInstance()
	runtimeContext := CreateRuntimeContext(system, stack)

	var actions []Action = make([]Action, 3)
	actions[0] = &dupOp
	actions[1] = &VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(1)}
	actions[2] = &addOp

	loop := CreateStartNextLoopAction(actions)

	err := loop.Apply(runtimeContext)
	if assert.NoError(t, err, "Loop should work") {
		if assert.Equal(t, 4, stack.Size()) {
			for i := 1; i <= 3; i++ {
				value, err := stack.Get(4 - i)
				if assert.NoError(t, err) {
					assert.Equal(t, int64(i), value.asNumericVar().value.IntPart())
				}
			}
		}
	}
}

func TestVariableDeclarationProgram(t *testing.T) {
	var v1 = CreateNumericVariableFromInt(4)
	var v2 = CreateNumericVariableFromInt(2)

	stack := CreateStack()
	stack.Push(v1)
	stack.Push(v2)

	system := CreateSystemInstance()
	runtimeContext := CreateRuntimeContext(system, stack)

	programVariable := CreateProgramVariable([]Action{
		&VariableEvaluationActionDesc{varName: "i"},
		&VariableEvaluationActionDesc{varName: "j"},
		&divOp, // use division such that is not commutative and this checks the order or arguments
	})

	varDecl := &VariableDeclarationActionDesc{
		varNames:           []string{"i", "j"},
		variableToEvaluate: programVariable,
	}

	err := varDecl.Apply(runtimeContext)
	if assert.NoError(t, err, "Error while running action") {
		if assert.Equal(t, 1, stack.Size()) {
			value, err := stack.PeekN(1)
			if assert.NoError(t, err, "No element on stack") {
				assert.Equal(t, int64(2), value[0].asNumericVar().value.IntPart())
			}
		}
	}
}
