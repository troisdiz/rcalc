package rcalc

import "fmt"

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
