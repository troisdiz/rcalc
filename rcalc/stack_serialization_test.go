package rcalc

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

func TestSaveAndReadStack(t *testing.T) {
	var stack Stack = CreateStack()
	v1 := CreateNumericVariableFromInt(2)
	v2 := CreateBooleanVariable(true)
	//v3 := CreateProgramVariable([]Action{})
	(&stack).Push(v1)
	(&stack).Push(v2)
	//(&stack).Push(v3)
	protoStack, err := CreateProtoFromStack(&stack)
	if assert.NoError(t, err) {
		out, err := proto.Marshal(protoStack)
		if assert.NoError(t, err) {
			readStack := protostack.Stack{}

			err = proto.Unmarshal(out, &readStack)
			if assert.NoError(t, err) {
				stack2, err := CreateStackFromProto(&readStack)
				if assert.NoError(t, err) {
					assert.Equal(t, &stack, stack2)
				}
			}
		}
	}

}
