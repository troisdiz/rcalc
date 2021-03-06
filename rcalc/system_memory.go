package rcalc

import "fmt"

type MemoryNode struct {
	parentFolder *MemoryNode
	name         string
}

func (node *MemoryNode) Name() string {
	return node.name
}

type MemoryVariable struct {
	MemoryNode
	value Variable
}

func (variable *MemoryVariable) Value() Variable {
	return variable.value
}

type MemoryFolder struct {
	MemoryNode
	subFolders []*MemoryFolder
	variables  []*MemoryVariable
}

func (folder *MemoryFolder) AsMemoryNode() *MemoryNode {
	return &folder.MemoryNode
}
func (folder *MemoryFolder) SubFolders() []*MemoryFolder {
	return folder.subFolders
}

func (folder *MemoryFolder) SubVariables() []*MemoryVariable {
	return folder.variables
}

type InternalMemory struct {
	memoryRoot    *MemoryFolder
	currentFolder *MemoryFolder
}

func (m *InternalMemory) getCurrentFolder() *MemoryFolder {
	return m.currentFolder
}

func (m *InternalMemory) getPath(node *MemoryNode) []string {
	var result []string

	for n := node; n.parentFolder != nil; n = n.parentFolder {
		result = append(result, n.Name())
	}

	resultLen := len(result)
	for i := 0; i < resultLen/2; i++ {
		result[i], result[resultLen-1-i] = result[resultLen-1-i], result[i]
	}
	return result
}

func (m *InternalMemory) resolvePath(path []string) *MemoryNode {
	//TODO implement me
	panic("implement me")
}

type Memory interface {
	getRoot() *MemoryFolder
	getCurrentFolder() *MemoryFolder
	//setCurrentFolder(f *MemoryFolder) error

	getPath(node *MemoryNode) []string
	resolvePath(path []string) *MemoryNode

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
		MemoryNode: MemoryNode{
			parentFolder: parent.AsMemoryNode(),
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
		MemoryNode: MemoryNode{name: variableName},
		value:      value,
	}
	parent.variables = append(parent.variables, memVar)
	return memVar, nil
}

func (im *InternalMemory) listVariables(parent *MemoryFolder) ([]*MemoryVariable, error) {
	return parent.variables[:], nil
}

func NewInternalMemory() *InternalMemory {
	homeFolder := &MemoryFolder{
		MemoryNode: MemoryNode{
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
