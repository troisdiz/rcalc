package rcalc

// System Access to non stack items : memory, exit function, etc
type System interface {
	exit()
	Memory() Memory
}

type SystemInternal interface {
	shouldStop() bool
}

type SystemInstance struct {
	shouldStopMarker bool
	memory           Memory
}

func (s *SystemInstance) shouldStop() bool {
	return s.shouldStopMarker
}

func (s *SystemInstance) exit() {
	s.shouldStopMarker = true
}

func (s *SystemInstance) Memory() Memory {
	return s.memory
}

func CreateSystemInstance() *SystemInstance {
	return &SystemInstance{
		shouldStopMarker: false,
		memory:           NewInternalMemory(),
	}
}

var EXIT_ACTION = NewOperationDesc(
	"quit",
	0,
	func(elts ...Variable) (bool, error) { return true, nil },
	0,
	func(system System, elts ...Variable) []Variable {
		system.exit()
		return nil
	})
