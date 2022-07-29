package rcalc

import (
	"bufio"
	"fmt"
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

		actions, parse_err := ParseToActions(cmds, "InteractiveShell", Registry)
		if parse_err != nil {
			message = parse_err.Error()
		} else {
			for _, action := range actions {
				if stack.Size() < action.NbArgs() {
					fmt.Printf("Not enough args on stack (%d vs %d)\n", stack.Size(), action.NbArgs())
					message = fmt.Sprintf("Not enough args on stack: only %d/%d available", action.NbArgs(), stack.Size())
					break
				} else {
					// TODO Handle error
					typesOK, _ := checkTypesForAction(&stack, action)
					if !typesOK {
						message = "Bad types on stack"
						break
					} else {
						applyErr := action.Apply(system, &stack)
						if applyErr != nil {
							message = applyErr.Error()
						}
					}
				}
				if system.shouldStop() {
					return
				}
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
