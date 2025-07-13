package rcalc

import (
	"fmt"

	"github.com/shopspring/decimal"
	"troisdizaines.com/rcalc/rcalc/protostack"
)

type ActionMarshallFunc func(reg *ActionRegistry, action Action) (*protostack.Action, error)
type ActionUnMarshallFunc func(reg *ActionRegistry, protoAction *protostack.Action) (Action, error)

type Action interface {
	OpCode() string
	NbArgs() int
	CheckTypes(elts ...Variable) (bool, error)
	Apply(runtimeContext *RuntimeContext) error
	Display() string
	MarshallFunc() ActionMarshallFunc
	UnMarshallFunc() ActionUnMarshallFunc
}

type ActionCommonDesc struct {
	opCode string
}

func (op *ActionCommonDesc) OpCode() string {
	return op.opCode
}
func ActionCommonDescMarshalFunction(reg *ActionRegistry, action Action) (*protostack.Action, error) {
	return &protostack.Action{
		Type:   protostack.ActionType_OPERATION,
		OpCode: action.OpCode(),
	}, nil
}

func (op *ActionCommonDesc) MarshallFunc() ActionMarshallFunc {
	return ActionCommonDescMarshalFunction
}

func ActionCommonDescUnMarshallFunc(reg *ActionRegistry, protoAction *protostack.Action) (Action, error) {
	protoOpCode := protoAction.OpCode
	if reg.ContainsOpCode(protoOpCode) {
		return reg.GetAction(protoOpCode), nil
	} else {
		return nil, fmt.Errorf("cannot find action with opcode %s", protoOpCode)
	}
}

func (op *ActionCommonDesc) UnMarshallFunc() ActionUnMarshallFunc {
	return ActionCommonDescUnMarshallFunc
}

type ActionApplyFn func(system System, stack *Stack) error

// ActionDesc implementation of Action interface
type ActionDesc struct {
	ActionCommonDesc
	nbArgs        int
	actionApplyFn ActionApplyFn
}

func (a *ActionDesc) MarshallFunc() ActionMarshallFunc {
	//TODO implement me
	panic("implement me ActionDesc MarshallFunc")
}

func (a *ActionDesc) UnMarshallFunc() ActionUnMarshallFunc {
	//TODO implement me
	panic("implement me ActionDesc UnMarshallFunc")
}

var _ Action = (*ActionDesc)(nil)

func NewActionDesc(opCode string, nbArgs int, checkTypeFn CheckTypeFn, applyFn ActionApplyFn) ActionDesc {
	return ActionDesc{
		ActionCommonDesc: ActionCommonDesc{
			opCode: opCode,
		},
		nbArgs:        nbArgs,
		actionApplyFn: applyFn,
	}
}
func (a *ActionDesc) NbArgs() int {
	return a.nbArgs
}

func (a *ActionDesc) CheckTypes(elts ...Variable) (bool, error) {
	return CheckNoop(elts...)
}

func (a *ActionDesc) Apply(runtimeContext *RuntimeContext) error {
	return a.actionApplyFn(runtimeContext.system, runtimeContext.stack)
}

func (a *ActionDesc) Display() string {
	return a.OpCode()
}

// OperationCommonDesc implementation of Action interface
type OperationCommonDesc struct {
	ActionCommonDesc
	nbArgs    int
	nbResults int
}

type CheckTypeFn func(elts ...Variable) (bool, error)
type OperationApplyFn func(system System, elts ...Variable) []Variable

type OperationDesc struct {
	OperationCommonDesc
	expandable  bool
	checkTypeFn CheckTypeFn
	applyFn     OperationApplyFn
}

var _ Action = (*OperationDesc)(nil)

func newOperationDesc(opCode string, nbArgs int, checkTypeFn CheckTypeFn, nbResults int, applyFn OperationApplyFn, expandable bool) OperationDesc {
	return OperationDesc{
		OperationCommonDesc: OperationCommonDesc{
			ActionCommonDesc: ActionCommonDesc{
				opCode: opCode,
			},
			nbArgs:    nbArgs,
			nbResults: nbResults,
		},
		expandable:  expandable,
		checkTypeFn: checkTypeFn,
		applyFn:     applyFn,
	}
}

func NewOperationDesc(opCode string, nbArgs int, checkTypeFn CheckTypeFn, nbResults int, applyFn OperationApplyFn) OperationDesc {
	return newOperationDesc(opCode, nbArgs, checkTypeFn, nbResults, applyFn, false)
}

func NewExpandableOperationDesc(opCode string, nbArgs int, checkTypeFn CheckTypeFn, nbResults int, applyFn OperationApplyFn) OperationDesc {
	return newOperationDesc(opCode, nbArgs, checkTypeFn, nbResults, applyFn, true)
}

func (op *OperationDesc) String() string {
	return fmt.Sprintf("Action(opCode = %s, nbArgs = %d)", op.opCode, op.nbArgs)
}

func (op *OperationDesc) NbArgs() int {
	return op.nbArgs
}

func (op *OperationDesc) CheckTypes(elts ...Variable) (bool, error) {

	// TODO This code has a lot of duplicate with Apply, next step is to change the
	// interface to make this create a context that Apply can use
	if op.expandable {
		argIsList := make([]bool, op.NbArgs())
		isExpanded := false
		listLength := 0
		for argIdx := 0; argIdx < op.NbArgs(); argIdx++ {
			if elts[argIdx].getType() == TYPE_LIST {
				argIsList[argIdx] = true
				listLength = elts[argIdx].asListVar().Size()
				isExpanded = true
			}
		}
		// TODO check on listLength consistency
		if isExpanded {
			for i := 0; i < listLength; i++ {
				tempInputs := make([]Variable, op.NbArgs())
				for inputIdx := 0; inputIdx < op.NbArgs(); inputIdx++ {
					if argIsList[inputIdx] {
						tempInputs[inputIdx] = elts[inputIdx].asListVar().items[i]
					} else {
						tempInputs[inputIdx] = elts[inputIdx]
					}
				}
				check, err := op.CheckTypes(tempInputs...)
				if err != nil {
					// TODO more precise error location message
					return check, err
				}
			}
			return true, nil
		} else {
			// not expanded case
			return op.checkTypeFn(elts...)
		}
	} else {
		return op.checkTypeFn(elts...)
	}
}

func (op *OperationDesc) NbResults() int {
	return op.nbResults
}

func (op *OperationDesc) Apply(runtimeContext *RuntimeContext) error {
	inputs, err := runtimeContext.stack.PopN(op.NbArgs())
	if err != nil {
		return err
	}
	if op.expandable {
		argIsList := make([]bool, op.NbArgs())
		isExpanded := false
		listLength := 0
		for argIdx := 0; argIdx < op.NbArgs(); argIdx++ {
			if inputs[argIdx].getType() == TYPE_LIST {
				argIsList[argIdx] = true
				listLength = inputs[argIdx].asListVar().Size()
				isExpanded = true
			}
		}
		// TODO check on listLength consistency
		if isExpanded {
			// We need as many result list as the applied function return results
			results := make([][]Variable, op.NbResults())
			for idx := range results {
				results[idx] = make([]Variable, listLength)
			}
			for i := 0; i < listLength; i++ {
				tempInputs := make([]Variable, op.NbArgs())
				for inputIdx := 0; inputIdx < op.NbArgs(); inputIdx++ {
					if argIsList[inputIdx] {
						tempInputs[inputIdx] = inputs[inputIdx].asListVar().items[i]
					} else {
						tempInputs[inputIdx] = inputs[inputIdx]
					}
				}
				tempResults := op.applyFn(runtimeContext.system, tempInputs...)

				for tmpResultIdx, tempResult := range tempResults {
					results[tmpResultIdx][i] = tempResult
				}
			}
			for _, result := range results {
				resultAsList := CreateListVariable(result)
				fmt.Printf("Result: %s\n", resultAsList.display())
				runtimeContext.stack.Push(resultAsList)
			}
		} else {
			// not expanded case
			results := op.applyFn(runtimeContext.system, inputs...)
			for _, elt := range results {
				runtimeContext.stack.Push(elt)
			}
		}
	} else {
		results := op.applyFn(runtimeContext.system, inputs...)
		for _, elt := range results {
			runtimeContext.stack.Push(elt)
		}
	}
	return nil
}

func (op *OperationDesc) Display() string {
	return op.opCode
}

type PureOperationApplyFn func(elts ...Variable) []Variable

func OpToActionFn(opFn PureOperationApplyFn) OperationApplyFn {
	return func(system System, elts ...Variable) []Variable {
		return opFn(elts...)
	}
}

type AlgebraicFn func(args ...decimal.Decimal) decimal.Decimal

type AlgebraicFunctionDesc struct {
	name      string
	argsCount int
	fn        AlgebraicFn
}

/* Registry stuff */

type ActionRegistry struct {
	actionDescs    map[string]Action
	dynamicActions map[string]struct {
		marshalFunc   ActionMarshallFunc
		unMarshalFunc ActionUnMarshallFunc
	}
	algebraicFunctionsByName map[string]AlgebraicFunctionDesc
}

func (reg *ActionRegistry) Register(aDesc Action) {
	reg.actionDescs[aDesc.OpCode()] = aDesc
}

type ActionPackage struct {
	staticActions       []Action
	dynamicActions      []Action
	algrebraicFunctions []AlgebraicFunctionDesc
}

func (ap *ActionPackage) AddStatic(action Action) {
	ap.staticActions = append(ap.staticActions, action)
}

func (ap *ActionPackage) AddDynamic(action Action) {
	ap.dynamicActions = append(ap.dynamicActions, action)
}

func (ap *ActionPackage) AddAlgebraicFunction(desc AlgebraicFunctionDesc) {
	ap.algrebraicFunctions = append(ap.algrebraicFunctions, desc)
}

func (reg *ActionRegistry) RegisterActions(aPackage *ActionPackage) {
	for _, aDesc := range aPackage.staticActions {
		reg.actionDescs[aDesc.OpCode()] = aDesc
	}
	for _, dynAction := range aPackage.dynamicActions {
		reg.dynamicActions[dynAction.OpCode()] = struct {
			marshalFunc   ActionMarshallFunc
			unMarshalFunc ActionUnMarshallFunc
		}{
			marshalFunc:   dynAction.MarshallFunc(),
			unMarshalFunc: dynAction.UnMarshallFunc(),
		}
	}
	for _, algFnDesc := range aPackage.algrebraicFunctions {
		reg.algebraicFunctionsByName[algFnDesc.name] = algFnDesc
	}
}

func initRegistry() *ActionRegistry {
	reg := ActionRegistry{
		actionDescs: map[string]Action{},
		dynamicActions: map[string]struct {
			marshalFunc   ActionMarshallFunc
			unMarshalFunc ActionUnMarshallFunc
		}{},
		algebraicFunctionsByName: map[string]AlgebraicFunctionDesc{},
	}
	reg.RegisterActions(&ArithmeticPackage)
	reg.RegisterActions(&TrigonometricPackage)
	reg.RegisterActions(&BooleanLogicPackage)
	reg.RegisterActions(&StatPackage)
	reg.RegisterActions(&StackPackage)
	reg.RegisterActions(&MemoryPackage)
	reg.RegisterActions(&StructOpsPackage)
	reg.RegisterActions(&ListPackage)
	reg.Register(&DebugOp)
	reg.Register(&VersionOp)
	reg.Register(&EXIT_ACTION)
	return &reg
}

func (reg *ActionRegistry) ContainsOpCode(opCode string) bool {
	_, ok := reg.actionDescs[opCode]
	return ok
}

func (reg *ActionRegistry) GetAction(opCode string) Action {
	actionDesc, ok := reg.actionDescs[opCode]
	if !ok {
		return nil
	} else {
		return actionDesc
	}
}

func (reg *ActionRegistry) GetDynamicActionMarshallFunc(opCode string) ActionMarshallFunc {
	if dynAction, ok := reg.dynamicActions[opCode]; ok {
		return dynAction.marshalFunc
	} else {
		return nil
	}
}

func (reg *ActionRegistry) GetDynamicActionUnMarshallFunc(opCode string) ActionUnMarshallFunc {
	if dynAction, ok := reg.dynamicActions[opCode]; ok {
		return dynAction.unMarshalFunc
	} else {
		return nil
	}
}

func (reg *ActionRegistry) GetAlgebraicFunction(fnName string) AlgebraicFn {
	algebraicFunctionDesc, ok := reg.algebraicFunctionsByName[fnName]
	if !ok {
		return nil
	} else {
		return algebraicFunctionDesc.fn
	}
}

func (reg *ActionRegistry) CreateActionFromProto(protoAction *protostack.Action) (Action, error) {
	if mFuncs, ok := reg.dynamicActions[protoAction.OpCode]; ok {
		action, err := mFuncs.unMarshalFunc(reg, protoAction)
		if err != nil {
			return nil, err
		}
		return action, nil
	} else {
		action := reg.GetAction(protoAction.OpCode)
		if action != nil {
			return action, nil
		}
	}
	return nil, fmt.Errorf("no unMarshallFunction found for type %d / OpCode %s", protoAction.Type, protoAction.OpCode)
}

func (reg *ActionRegistry) GetDynamicActionOpCodes() []string {
	result := make([]string, len(reg.dynamicActions))
	for k := range reg.dynamicActions {
		result = append(result, k)
	}
	return result
}

var Registry = initRegistry()
