package rcalc

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

func getDynamicActionsForTesting() []Action {
	return []Action{
		&VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(1)},
		&StartNextLoopActionDesc{
			actions: []Action{
				&VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(7)},
			},
		},
		&ForNextLoopActionDesc{
			varName: "n",
			actions: []Action{&VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(7)}},
		},
		&IfThenElseActionDesc{
			ifActions:   []Action{&VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(7)}},
			thenActions: []Action{&VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(8)}},
			elseActions: nil,
		},
		&VariableDeclarationActionDesc{
			varNames: []string{"n", "m"},
			programVariable: CreateProgramVariable([]Action{&VariablePutOnStackActionDesc{
				value: CreateNumericVariableFromInt(7)}}),
		},
		&EvalProgramActionDesc{
			program: CreateProgramVariable([]Action{&VariablePutOnStackActionDesc{
				value: CreateNumericVariableFromInt(7)}}),
		},
		&VariableEvaluationActionDesc{
			varName: "a",
		},
	}
}

/**
Test that all registered dynamic actions are tested
*/
func TestAllDynamicActionsAreTested(t *testing.T) {
	var dynamicActionOpCodes = make(map[string]interface{})
	for dynAction := range Registry.dynamicActions {
		dynamicActionOpCodes[dynAction] = nil
	}
	for _, action := range getDynamicActionsForTesting() {
		delete(dynamicActionOpCodes, action.OpCode())
	}
	assert.Len(t, dynamicActionOpCodes, 0)
}

func TestSaveAndReadActions(t *testing.T) {

	for _, dynAction := range getDynamicActionsForTesting() {
		t.Run(dynAction.OpCode(), func(t *testing.T) {
			protoAction, err := Registry.GetDynamicActionMarshallFunc(dynAction.OpCode())(Registry, dynAction)
			if assert.NoError(t, err) {
				serializedProto, err := proto.Marshal(protoAction)
				if assert.NoError(t, err) {
					unSerializedProtoAction := &protostack.Action{}
					err = proto.Unmarshal(serializedProto, unSerializedProtoAction)

					if assert.NoError(t, err) {
						fromProto, err := Registry.CreateActionFromProto(unSerializedProtoAction)
						if assert.NoError(t, err) {
							assert.Equal(t, dynAction, fromProto)
						}
					}
				}
			}
		})
	}
}

func TestSaveAndReadStack(t *testing.T) {

	var stack *Stack = CreateStack()
	v1 := CreateNumericVariableFromInt(2)
	stack.Push(v1)

	v2 := CreateBooleanVariable(true)
	stack.Push(v2)

	a1 := &VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(5)}
	a2 := &VariablePutOnStackActionDesc{value: CreateNumericVariableFromInt(7)}
	v3 := CreateProgramVariable([]Action{a1, a2, &addOp})
	stack.Push(v3)

	protoStack, err := CreateProtoFromStack(stack)
	if assert.NoError(t, err) {
		out, err := proto.Marshal(protoStack)
		if assert.NoError(t, err) {
			readStack := protostack.Stack{}

			err = proto.Unmarshal(out, &readStack)
			if assert.NoError(t, err) {
				stack2, err := CreateStackFromProto(Registry, &readStack)
				if assert.NoError(t, err) {
					assert.Equal(t, stack, stack2)
				}
			}
		}
	}

}
