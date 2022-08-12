package rcalc

import "fmt"

type RuntimeContext struct {
	system    System
	stack     *Stack
	rootScope *Scope
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
	// TODO memory variables
	return nil, fmt.Errorf("variable named %s not found", varName)
}

func (s *Scope) SetVariableValue(varName string, value Variable) error {
	// TODO: should we check for shadowing ?
	s.variables[varName] = value
	return nil
}
