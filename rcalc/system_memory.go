package rcalc

import "fmt"

type MemoryNode interface {
	asMemoryVariable() *MemoryVariable
	asMemoryFolder() *MemoryFolder
	getParent() *MemoryFolder
	Name() string
}

type AbstractMemoryNode struct {
	parentFolder *MemoryFolder
	name         string
}

var _ MemoryNode = (*AbstractMemoryNode)(nil)

func (node *AbstractMemoryNode) getParent() *MemoryFolder {
	return node.parentFolder
}

func (node *AbstractMemoryNode) Name() string {
	return node.name
}

func (node *AbstractMemoryNode) asMemoryVariable() *MemoryVariable {
	panic("Cannot cast")
}

func (node *AbstractMemoryNode) asMemoryFolder() *MemoryFolder {
	panic("Cannot cast")
}

type MemoryVariable struct {
	AbstractMemoryNode
	value Variable
}

var _ MemoryNode = (*MemoryVariable)(nil)

func (mv *MemoryVariable) asMemoryVariable() *MemoryVariable {
	return mv
}

func (mv *MemoryVariable) Value() Variable {
	return mv.value
}

type MemoryFolder struct {
	AbstractMemoryNode
	subFolders []*MemoryFolder
	variables  []*MemoryVariable
}

var _ MemoryNode = (*MemoryFolder)(nil)

func (folder *MemoryFolder) asMemoryFolder() *MemoryFolder {
	return folder
}

func (folder *MemoryFolder) SubFolders() []*MemoryFolder {
	return folder.subFolders
}

func (folder *MemoryFolder) SubVariables() []*MemoryVariable {
	return folder.variables
}

func (folder *MemoryFolder) SubNodes() []MemoryNode {
	var result []MemoryNode
	for _, elt := range folder.SubFolders() {
		result = append(result, elt)
	}
	for _, elt := range folder.SubVariables() {
		result = append(result, elt)
	}
	return result
}

type InternalMemory struct {
	memoryRoot    *MemoryFolder
	currentFolder *MemoryFolder
}

func (m *InternalMemory) getCurrentFolder() *MemoryFolder {
	return m.currentFolder
}

func (m *InternalMemory) getPath(node MemoryNode) []string {
	var result []string

	for n := node; n.getParent() != nil; n = n.getParent() {
		result = append(result, n.Name())
	}

	resultLen := len(result)
	for i := 0; i < resultLen/2; i++ {
		result[i], result[resultLen-1-i] = result[resultLen-1-i], result[i]
	}
	return result
}

func (m *InternalMemory) resolvePath(path []string) MemoryNode {
	pathNode := m.getRoot()
	totalDepth := len(path)
	if totalDepth == 0 {
		return pathNode
	}
	for _, pathElt := range path[:totalDepth-1] {
		var nextFolder *MemoryFolder
		for _, subFolder := range pathNode.asMemoryFolder().subFolders {
			if subFolder.name == pathElt {
				nextFolder = subFolder
				break
			}
		}
		if nextFolder == nil {
			return nil
		} else {
			pathNode = nextFolder
		}
	}
	for _, subVariable := range pathNode.SubNodes() {
		if subVariable.Name() == path[totalDepth-1] {
			return subVariable
		}
	}
	return nil
}

type Memory interface {
	getRoot() *MemoryFolder
	getCurrentFolder() *MemoryFolder
	//setCurrentFolder(f *MemoryFolder) error

	getPath(node MemoryNode) []string
	resolvePath(path []string) MemoryNode

	createFolder(folderName string, parent *MemoryFolder) (*MemoryFolder, error)
	createVariable(variableName string, parent *MemoryFolder, value Variable) (*MemoryVariable, error)
	listVariables(parent *MemoryFolder) ([]*MemoryVariable, error)
	/*
		cd(path string)
		currentDir() string
		list(path string) []string
	*/
}

func (im *InternalMemory) getRoot() *MemoryFolder {
	return im.memoryRoot
}

func (im *InternalMemory) createFolder(folderName string, parent *MemoryFolder) (*MemoryFolder, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent folder is nil")
	}
	newFolder := &MemoryFolder{
		AbstractMemoryNode: AbstractMemoryNode{
			parentFolder: parent,
			name:         folderName,
		},
	}
	parent.subFolders = append(parent.SubFolders(), newFolder)
	return newFolder, nil
}

func (im *InternalMemory) createVariable(variableName string, parent *MemoryFolder, value Variable) (*MemoryVariable, error) {
	if parent == nil {
		return nil, fmt.Errorf("Cannot create memory variable with nil parent folder")
	}
	memVar := &MemoryVariable{
		AbstractMemoryNode: AbstractMemoryNode{name: variableName},
		value:              value,
	}
	parent.variables = append(parent.variables, memVar)
	return memVar, nil
}

func (im *InternalMemory) listVariables(parent *MemoryFolder) ([]*MemoryVariable, error) {
	return parent.variables[:], nil
}

func NewInternalMemory() *InternalMemory {
	homeFolder := &MemoryFolder{
		AbstractMemoryNode: AbstractMemoryNode{
			name:         "HOME",
			parentFolder: nil},
		subFolders: nil,
		variables:  nil,
	}
	return &InternalMemory{
		memoryRoot:    homeFolder,
		currentFolder: homeFolder,
	}
}
