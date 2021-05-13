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

func findNearestFork() *Node {
	if Current == nil {
		return nil
	}
	var node *Node
	for node = Current; node.NodeTypeLabel != Fork && node != nil; node = node.Parent {
	}
	return node
}
