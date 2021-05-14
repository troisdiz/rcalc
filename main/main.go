package main

import (
	"bufio"
	"github.com/troisdiz/rcalc/rcalc"
	"os"
	"strings"
)

func main() {

	var stack rcalc.Stack = rcalc.Create()

	for {
		// print stack
		rcalc.DisplayStack(stack, 3)

		// print prompt
		input := bufio.NewScanner(os.Stdin)

		// wait for cmd
		input.Scan()

		// interpret cmd
		var cmds = input.Text()

		singleCmd := strings.TrimSpace(cmds)
		if singleCmd == "quit" {
			return
		}

	}



}
