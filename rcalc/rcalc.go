package rcalc

import (
	"bufio"
	"fmt"
	"os"
)

func Run() {

	var stack Stack = CreateStack()
	var message string = ""
	var system *SystemInstance = CreateSystemInstance()
	for {
		// print stack
		DisplayStack(stack, message, 3)

		// print prompt
		input := bufio.NewScanner(os.Stdin)

		// wait for cmd
		input.Scan()

		// interpret cmd
		var cmds = input.Text()

		var expressions []*ExprElement
		expressions, _ = ParseExpression(Registry, cmds)

		for _, expr := range expressions {
			switch expr.eltType {
			case ACTION_EXPR_TYPE:
				action := expr.asAction()
				var stackElts []StackElt = make([]StackElt, action.NbArgs())
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
							for i := 0; i < action.NbArgs(); i++ {
								stackElt, err := stack.Pop()
								if err != nil {
									panic("Stack error !!")
								}
								stackElts[i] = stackElt
							}
							stackEltResult := action.Apply(system, stackElts...)
							if stackEltResult != nil {
								for _, stackElt := range stackEltResult {
									stack.Push(stackElt)
								}
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

func checkTypesForAction(s *Stack, a Action) (bool, error) {
	elts, err := s.Peek(a.NbArgs())
	ok, err := a.CheckTypes(elts...)
	if err != nil {
		return false, err
	} else {
		return ok, nil
	}
}
