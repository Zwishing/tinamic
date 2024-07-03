package model

import (
	"github.com/gofrs/uuid"
	"time"
)

// Node is the interface for all nodes
type StoreNode interface {
	GetTitle() string
	GetKey() string
	GetType() string
	GetSize() int64
	GetModifiedTime() time.Time
}

// FolderNode represents a folder node that can have children
type FolderNode struct {
	Title        string      `json:"title"`
	Key          string      `json:"key"`
	Type         string      `json:"type"`
	Size         int64       `json:"size"`
	ModifiedTime time.Time   `json:"modifiedTime"`
	Children     []StoreNode `json:"children"`
}

// FileNode represents a file node that cannot have children
type FileNode struct {
	Title        string    `json:"title"`
	Key          string    `json:"key"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modifiedTime"`
}

// FolderFileTree represents the root of the folder-file tree
type FolderFileTree struct {
	Root *FolderNode
}

func NewFolderNode(title string, size int64, modifiedTime time.Time) *FolderNode {
	key, _ := uuid.NewV4()
	return &FolderNode{
		Title:        title,
		Key:          key.String(),
		Type:         "folder",
		Size:         size,
		ModifiedTime: modifiedTime,
		Children:     nil,
	}
}

// GetTitle returns the title of the folder node
func (f *FolderNode) GetTitle() string {
	return f.Title
}

// GetKey returns the key of the folder node
func (f *FolderNode) GetKey() string {
	return f.Key
}

// GetType returns the type of the folder node
func (f *FolderNode) GetType() string {
	return "folder"
}

// GetSize returns the size of the folder node
func (f *FolderNode) GetSize() int64 {
	return f.Size
}

// GetModifiedTime returns the modified time of the folder node
func (f *FolderNode) GetModifiedTime() time.Time {
	return f.ModifiedTime
}

// AddChild adds a child node to the folder node
func (f *FolderNode) AddChild(child StoreNode) {
	f.Children = append(f.Children, child)
}

// GetTitle returns the title of the file node
func (file *FileNode) GetTitle() string {
	return file.Title
}

// GetKey returns the key of the file node
func (file *FileNode) GetKey() string {
	return file.Key
}

// GetType returns the type of the file node
func (file *FileNode) GetType() string {
	return "file"
}

// GetSize returns the size of the file node
func (file *FileNode) GetSize() int64 {
	return file.Size
}

// GetModifiedTime returns the modified time of the file node
func (file *FileNode) GetModifiedTime() time.Time {
	return file.ModifiedTime
}

// AddNode adds a new node to the tree at the specified parent key
func (tree *FolderFileTree) AddNode(parentKey string, newNode StoreNode) bool {
	parentNode := FindNode(tree.Root, parentKey)
	if parentNode != nil {
		if folder, ok := parentNode.(*FolderNode); ok {
			folder.AddChild(newNode)
			return true
		}
	}
	return false
}

// FindNode finds a node by key
func FindNode(node StoreNode, key string) StoreNode {
	if node.GetKey() == key {
		return node
	}
	if folder, ok := node.(*FolderNode); ok {
		for _, child := range folder.Children {
			if found := FindNode(child, key); found != nil {
				return found
			}
		}
	}
	return nil
}
