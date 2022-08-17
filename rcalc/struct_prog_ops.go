package rcalc

import "fmt"

type VariablePutOnStackActionDesc struct {
	value Variable
}

func (a *VariablePutOnStackActionDesc) NbArgs() int {
	return 0
}

func (a *VariablePutOnStackActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (a *VariablePutOnStackActionDesc) Apply(runtimeContext *RuntimeContext) error {
	runtimeContext.stack.Push(a.value)
	return nil
}

func (a *VariablePutOnStackActionDesc) OpCode() string {
	return "__hidden__" + "PutOnStack"
}

func (a *VariablePutOnStackActionDesc) String() string {
	return fmt.Sprintf("%s(%s)", a.OpCode(), a.value.String())
}

type VariableEvaluationActionDesc struct {
	varName string
}

func (a *VariableEvaluationActionDesc) NbArgs() int {
	return 0
}

func (a *VariableEvaluationActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return true, nil
}

func (a *VariableEvaluationActionDesc) Apply(runtimeContext *RuntimeContext) error {
	value, err := runtimeContext.GetVariableValue(a.varName)
	if err != nil {
		return err
	}
	runtimeContext.stack.Push(value)
	return nil
}

func (a *VariableEvaluationActionDesc) OpCode() string {
	return "__hidden__" + "VariableEvaluation"
}

func (a *VariableEvaluationActionDesc) String() string {
	return fmt.Sprintf("%s(%s)", a.OpCode(), a.varName)
}

type StartNextLoopActionDesc struct {
	actions []Action
}

func (a *StartNextLoopActionDesc) NbArgs() int {
	return 2
}

func (a *StartNextLoopActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	for i := 0; i <= 1; i++ {
		if elts[i].getType() != TYPE_NUMERIC && elts[i].asNumericVar().value.IsInteger() {
			return false, fmt.Errorf("%s at stack level %d is not an integer", elts[i].String(), i+1)
		}
	}
	return true, nil
}

func (a *StartNextLoopActionDesc) Apply(runtimeContext *RuntimeContext) error {
	boundaries, err := runtimeContext.stack.PopN(2)
	if err != nil {
		return err
	}
	start := boundaries[0].asNumericVar().value.IntPart()
	end := boundaries[1].asNumericVar().value.IntPart()
	for i := start; i <= end; i++ {
		for _, action := range a.actions {
			err = action.Apply(runtimeContext)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *StartNextLoopActionDesc) OpCode() string {
	return "__hidden__" + "StartNextLoop"
}

func (a *StartNextLoopActionDesc) String() string {
	return fmt.Sprintf("%s ()", a.OpCode())
}

func CreateStartNextLoopAction(actions []Action) *StartNextLoopActionDesc {
	return &StartNextLoopActionDesc{actions: actions}
}

type ForNextLoopActionDesc struct {
	varName string
	actions []Action
}

func (a ForNextLoopActionDesc) OpCode() string {
	return "__hidden__" + "ForNextLoop"
}

func (a ForNextLoopActionDesc) NbArgs() int {
	return 2
}

func (a ForNextLoopActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	for i := 0; i <= 1; i++ {
		if elts[i].getType() != TYPE_NUMERIC && elts[i].asNumericVar().value.IsInteger() {
			return false, fmt.Errorf("%s at stack level %d is not an integer", elts[i].String(), i+1)
		}
	}
	return true, nil
}

func (a ForNextLoopActionDesc) Apply(runtimeContext *RuntimeContext) error {

	runtimeContext.EnterNewScope()
	defer func() {
		runtimeContext.LeaveScope()
	}()

	boundaries, err := runtimeContext.stack.PopN(2)
	if err != nil {
		return err
	}
	start := boundaries[0].asNumericVar().value.IntPart()
	end := boundaries[1].asNumericVar().value.IntPart()
	// new scope

	for i := start; i <= end; i++ {
		// set var value
		err := runtimeContext.SetVariableValue(a.varName, CreateNumericVariableFromInt(int(i)))
		if err != nil {
			return err
		}
		for _, action := range a.actions {
			err = action.Apply(runtimeContext)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
