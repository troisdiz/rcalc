package rcalc

import "fmt"

type RuntimeContext struct {
	system       System
	stack        *Stack
	rootScope    *Scope
	currentScope *Scope
}

type VariableReader interface {
	GetVariableValue(varName string) (Variable, error)
}

func CreateRuntimeContext(system System, stack *Stack) *RuntimeContext {
	rtContext := &RuntimeContext{
		system: system,
		stack:  stack,
		rootScope: &Scope{
			parent:    nil,
			rt:        nil,
			variables: make(map[string]Variable),
		},
	}
	rtContext.rootScope.rt = rtContext
	return rtContext
}

func (rt *RuntimeContext) RunAction(action Action) error {
	if rt.stack.Size() < action.NbArgs() {
		// fmt.Printf("Not enough args on stack (%d vs %d)\n", rt.stack.Size(), action.NbArgs())
		return fmt.Errorf("not enough args on stack: only %d/%d available", action.NbArgs(), rt.stack.Size())
	} else {
		// TODO Handle error
		typesOK, err := checkTypesForAction(rt.stack, action)
		if !typesOK {
			return err
		} else {
			applyErr := action.Apply(rt)
			if applyErr != nil {
				return applyErr
			}
		}
	}
	return nil
}

type Scope struct {
	parent    *Scope
	rt        *RuntimeContext
	variables map[string]Variable
}

func (s *Scope) GetVariableValue(varName string) (Variable, error) {
	lookupScope := s
	for lookupScope != nil {
		if val, ok := lookupScope.variables[varName]; ok {
			return val, nil
		}
		lookupScope = lookupScope.parent
	}
	// Nothing found in local vairables
	return nil, fmt.Errorf("variable named %s not found", varName)
}

func (s *Scope) SetVariableValue(varName string, value Variable) error {
	// TODO: should we check for shadowing ?
	s.variables[varName] = value
	return nil
}

func (rt *RuntimeContext) GetVariableValue(varName string) (Variable, error) {

	value, err := rt.currentScope.GetVariableValue(varName)
	// Let's look in main Memory
	if err != nil {
		memory := rt.system.Memory()
		currentFolder := memory.getCurrentFolder()
		varPath := append(memory.getPath(currentFolder), varName)
		node := memory.resolvePath(varPath)
		varNode := node.asMemoryVariable()
		value = varNode.value
		err = nil
	}
	return value, err
}

func (rt *RuntimeContext) SetVariableValue(varName string, value Variable) error {
	return rt.currentScope.SetVariableValue(varName, value)
}

func (rt *RuntimeContext) EnterNewScope() {
	newScope := &Scope{
		parent:    rt.currentScope,
		rt:        rt,
		variables: make(map[string]Variable),
	}
	rt.currentScope = newScope
}

func (rt *RuntimeContext) LeaveScope() {
	rt.currentScope.rt = nil
	rt.currentScope = rt.currentScope.parent
}
