package semanthoid

func CreateFork() {
	fork := &Node{NodeTypeLabel: Fork, Identifier: ForkIdentifier}
	if Root == nil {
		Root = fork
		Current = fork
	} else {
		Current.Right = fork
		fork.Parent = Current
		Current = fork
	}
}

func findNearestForkFromCurrent() *Node {
	if Current == nil {
		return nil
	}
	var node *Node
	for node = Current; node.NodeTypeLabel != Fork && node != nil; node = node.Parent {
	}
	return node
}

func findNearestForkFromRootAmongLeft() *Node {
	if Root == nil {
		return nil
	}
	node := Root
	for {
		if node.Right != nil && node.Right.NodeTypeLabel == Fork {
			return node.Right
		} else if node.Left != nil {
			node = node.Left
		} else {
			return nil
		}
	}
}
