package rcalc

var toListOp = NewRawStackOpWithCheck("tolist", 1, CheckFirstInt, func(system System, stack *Stack) error {
	n, err := stack.Pop()
	if err != nil {
		return err
	}
	stackElts, err := stack.PopN(int(n.asNumericVar().value.IntPart()))
	if err != nil {
		return err
	}
	listVar := CreateListVariable(stackElts)
	stack.Push(listVar)
	return nil
})

var ListPackage = ActionPackage{
	staticActions: []Action{
		&toListOp,
	},
	dynamicActions: []Action{},
}
