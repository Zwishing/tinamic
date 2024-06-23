package model

import (
	"fmt"
	"testing"
)

func TestFindNode(t *testing.T) {
	tree := FolderFileTree{
		Root: &FolderNode{Key: "1"},
	}
	tree.AddNode("1", &FolderNode{
		Key: "2",
	})
	fmt.Println(tree.Root.Children[0].GetKey())
}
