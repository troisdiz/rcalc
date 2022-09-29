package rcalc

import (
	"bufio"
	"os"
)

func Run() {

	var stack = CreateStack()
	var message = ""
	var system = CreateSystemInstance()
	for {
		// print stack
		DisplayStack(stack, message, 3, true)

		// print prompt
		input := bufio.NewScanner(os.Stdin)

		// wait for cmd
		input.Scan()

		// interpret cmd
		var cmds = input.Text()

		// Message to display above (temp way of doing this)
		message = ""

		actions, parseErr := ParseToActions(cmds, "InteractiveShell", Registry)
		if parseErr != nil {
			message = parseErr.Error()
		} else {
			runtimeContext := CreateRuntimeContext(system, stack)
			err := stack.StartSession()
			if err != nil {
				return
			}
			for _, action := range actions {
				err := runtimeContext.RunAction(action)
				if err != nil {
					message = err.Error()
					// in case of error, stop evaluation
					break
				} else {
					message = ""
				}
				if system.shouldStop() {
					err := stack.CloseSession()
					if err != nil {
						return
					}
					return
				}
			}
			err = stack.CloseSession()
			if err != nil {
				return
			}
		}
	}
}

func checkTypesForAction(s *Stack, a Action) (bool, error) {
	elts, _ := s.PeekN(a.NbArgs())
	ok, err := a.CheckTypes(elts...)
	if err != nil {
		return false, err
	} else {
		return ok, nil
	}
}
