package semanthoid

import (
	"cpp-custom/logger"
	"errors"
)

func ClearCurrentRightSubTree() error {
	rightSubTreeRoot := findRightSubTreeRoot()
	if rightSubTreeRoot == nil {
		return errors.New("fatal invalid cleanup attempt")
	}
	if rightSubTreeRoot.Parent == nil {
		Root = nil
		Current = nil
	} else {
		rightSubTreeRoot.Parent.Right = nil
	}
	logger.Log("memory_l", "right subtree release of memory")
	logger.Log("tree_l", "end of block and release of memory\n"+TreeToString())
	return nil
}

func findRightSubTreeRoot() *Node {
	if Current == nil {
		return nil
	}
	var node *Node
	for node = Current; node.Parent != nil && node.Parent.Right != node; node = node.Parent {
	}
	return node
}
