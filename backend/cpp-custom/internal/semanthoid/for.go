package semanthoid

func CreateForOperator(counterType int, counterIdentifier string, counterValue *DataTypeValue, forBody *Node) *Node {
	node := &Node{
		NodeTypeLabel: For,
		Identifier:    "for",
	}
	if counterType != Error {
		node.Right = &Node{
			NodeTypeLabel: Variable,
			Identifier:    counterIdentifier,
			DataTypeLabel: counterType,
			DataValue:     counterValue,
			Parent:        node,
		}
		node.Right.Left = forBody
		forBody.Parent = node.Right
	} else {
		node.Right = forBody
		forBody.Parent = node
	}
	return node
}
