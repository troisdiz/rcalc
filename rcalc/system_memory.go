package rcalc

import "fmt"

type MemoryNode struct {
	name string
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
	parentFolder *MemoryFolder
	subFolders   []*MemoryFolder
	variables    []*MemoryVariable
}

func (folder *MemoryFolder) SubFolders() []*MemoryFolder {
	return folder.subFolders
}

func (folder *MemoryFolder) SubVariables() []*MemoryVariable {
	return folder.variables
}

type InternalMemory struct {
	memoryRoot *MemoryFolder
}

type Memory interface {
	getRoot() *MemoryFolder

	createFolder(folderName string, parent *MemoryFolder) error
	createVariable(variableName string, parent *MemoryFolder, value Variable) error
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

func (im *InternalMemory) createFolder(folderName string, parent *MemoryFolder) error {
	//TODO implement me
	panic("implement me")
}

func (im *InternalMemory) createVariable(variableName string, parent *MemoryFolder, value Variable) error {
	if parent == nil {
		return fmt.Errorf("Cannot create memory variable with nil parent folder")
	}
	memVar := &MemoryVariable{
		MemoryNode: MemoryNode{name: variableName},
		value:      value,
	}
	parent.variables = append(parent.variables, memVar)
	return nil
}

func (im *InternalMemory) listVariables(parent *MemoryFolder) ([]*MemoryVariable, error) {
	return parent.variables[:], nil
}

func NewInternalMemory() *InternalMemory {
	return &InternalMemory{memoryRoot: &MemoryFolder{
		MemoryNode:   MemoryNode{name: "ROOT"},
		parentFolder: nil,
		subFolders:   nil,
		variables:    nil,
	}}
}
