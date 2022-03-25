package rcalc

/**
Access to non stack items : memory, exit function, etc
*/
type System interface {
	exit()
}

type SystemInternal interface {
	shouldStop() bool
}

type SystemInstance struct {
	shouldStopMarker bool
}

func (s *SystemInstance) shouldStop() bool {
	return s.shouldStopMarker
}

func (s *SystemInstance) exit() {
	s.shouldStopMarker = true
}

func CreateSystemInstance() *SystemInstance {
	return &SystemInstance{
		shouldStopMarker: false,
	}
}

var EXIT_ACTION = ActionDesc{
	opCode: "quit",
	nbArgs: 0,
	checkTypeFn: func(elts ...StackElt) (bool, error) {
		return true, nil
	},
	applyFn: func(system System, elts ...StackElt) []StackElt {
		system.exit()
		return nil
	},
}
