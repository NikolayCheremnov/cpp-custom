package semanthoid

func CreateCompositeOperator() *Node {
	if Current.NodeTypeLabel != CompositeOperator { // composite operator node does not exist
		Current.Left = &Node{
			NodeTypeLabel: CompositeOperator,
			Identifier:    CompositeOperatorIdentifier,
		}
	}
	Current = Current.Left
	return Current
}
