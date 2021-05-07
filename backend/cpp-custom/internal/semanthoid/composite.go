package semanthoid

func CreateCompositeOperator() *Node {
	return &Node{
		NodeTypeLabel: CompositeOperator,
		Identifier:    CompositeOperatorIdentifier,
	}
}
