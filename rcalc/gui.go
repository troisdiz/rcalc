package rcalc

import "fmt"

func DisplayStack(s Stack, minElts int) {
	stackSize := s.Size()
	for i := minElts-1; i >= stackSize; i-- {
		displayStackLevel(i, nil)
	}
	for i := stackSize -1; i >= 0; i-- {
		elt, _ := s.Get(i)
		displayStackLevel(i, elt)
	}
}

func displayStackLevel(level int, elt StackElt) {
	var value string = ""
	if elt != nil {
		value = elt.display()
	}
	fmt.Printf("%2d: %10s", level+1, value)
	fmt.Println()
}
