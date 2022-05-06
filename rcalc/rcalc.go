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

		var expressions []*ExprElement
		var err error

		message = ""

		expressions, err = ParseExpression(Registry, cmds)
		if err != nil {
			message = err.Error()
		} else {

			for _, expr := range expressions {
				switch expr.eltType {
				case ACTION_EXPR_TYPE:
					action := expr.asAction()
					if stack.Size() < action.NbArgs() {
						fmt.Printf("Not enough args on stack (%d vs %d)\n", stack.Size(), action.NbArgs())
					} else {
						typesOK, err := checkTypesForAction(&stack, action)
						if err != nil {
							panic(fmt.Sprintf("Error while checking types of %s : %v", action.OpCode(), err))
						} else {
							if !typesOK {
								message = "Bad types on stack"
							} else {
								applyErr := action.Apply(system, &stack)
								if applyErr != nil {
									message = applyErr.Error()
								}
							}
						}
					}
					if system.shouldStop() {
						return
					}
				case STACK_ELT_EXPR_TYPE:
					ste := expr.asStackElt()
					stack.Push(ste)
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
