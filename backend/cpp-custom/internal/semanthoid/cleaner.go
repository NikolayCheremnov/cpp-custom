package semanthoid

import (
	"cpp-custom/logger"
)

func ClearCurrentRightSubTree() error {
	rightSubTreeRoot := findNearestFork()
	if rightSubTreeRoot == nil {
		logger.Log("memory_l", "right subtree was empty, release of memory is not necessary")
	} else {
		if rightSubTreeRoot.Parent == nil {
			Root = nil
			Current = nil
		} else {
			rightSubTreeRoot.Parent.Right = nil
			Current = rightSubTreeRoot.Parent
		}
		logger.Log("memory_l", "right subtree release of memory")
	}
	logger.Log("tree_l", "end of block and release of memory\n"+TreeToString())
	return nil
}
