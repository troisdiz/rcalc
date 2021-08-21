package rcalc

import (
	"bufio"
	"os"
)

func Run() {

	var stack Stack = CreateStack()
	var system *SystemInstance = CreateSystemInstance()
	for {
		// print stack
		DisplayStack(stack, 3)

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
				for i := 0; i < action.NbArgs(); i++ {
					stackElt, err := stack.Pop()
					if err != nil {
						panic("Stack error !!")
					}
					stackElts[i] = stackElt
				}
				stackEltResult := action.Apply(system, stackElts...)
				if stackEltResult != nil {
					stack.Push(stackEltResult)
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
