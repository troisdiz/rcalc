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
	runtimeContext := CreateRuntimeContext(system, &stack)

	var actions []Action = make([]Action, 3)
	actions[0] = &dupOp
	actions[1] = &VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(1)}
	actions[2] = &addOp

	loop := CreateStartNextLoopAction(actions)

	err := loop.Apply(runtimeContext)
	if assert.NoError(t, err, "Loop should work") {
		if assert.Equal(t, stack.Size(), 4) {
			for i := 1; i <= 3; i++ {
				value, err := stack.Get(4 - i)
				if assert.NoError(t, err) {
					assert.Equal(t, int64(i), value.asNumericVar().value.IntPart())
				}
			}
		}
	}

}
