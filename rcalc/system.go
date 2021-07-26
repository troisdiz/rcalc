package rcalc

/**
Access to non stack items : memory, exit function, etc
 */
type System interface {
    exit()
}

var EXIT_ACTION = ActionDesc{
    opCode:      "quit",
    nbArgs:      0,
    checkTypeFn: func(elts ...StackElt) (bool, error) {
        return true, nil
    },
    applyFn: func(system System, elts ...StackElt) StackElt {
        return nil
    },
}
