package rcalc

type MemoryNode struct {
	name string
}

func (node *MemoryNode) Name() string {
	return node.name
}

type MemoryVariable struct {
	MemoryNode
	value StackElt
}

func (variable *MemoryVariable) Value() StackElt {
	return variable.value
}

type MemoryFolder struct {
	MemoryNode
	subFolders []MemoryFolder
	variables  []MemoryVariable
}

func (folder *MemoryFolder) SubFolders() []MemoryFolder {
	return folder.subFolders
}

func (folder *MemoryFolder) SubVariables() []MemoryVariable {
	return folder.variables
}

type InternalMemory struct {
	memoryRoot *MemoryFolder
}

func NewInternalMemory() *InternalMemory {
	return &InternalMemory{memoryRoot: &MemoryFolder{
		MemoryNode: MemoryNode{name: "ROOT"},
		subFolders: nil,
		variables:  nil,
	}}
}

// System Access to non stack items : memory, exit function, etc
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

var EXIT_ACTION = NewOperationDesc(
	"quit",
	0,
	func(elts ...StackElt) (bool, error) { return true, nil },
	func(system System, elts ...StackElt) []StackElt {
		system.exit()
		return nil
	})
