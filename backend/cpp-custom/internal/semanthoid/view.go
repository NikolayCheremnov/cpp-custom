package semanthoid

import (
	"strconv"
	"strings"
)

func (node *Node) ToString() string {
	res := ""
	if node == Current {
		res += "current: "
	}
	res += "node type " + strconv.Itoa(node.NodeTypeLabel) + ": identifier " + node.Identifier
	switch node.NodeTypeLabel { // special view for each node type
	case Variable, Constant:
		res += ": data type " + strconv.Itoa(node.DataTypeLabel)
		res += ": data value bool " + strconv.Itoa(node.DataValue.DataAsBool) + " int " + strconv.Itoa(node.DataValue.DataAsInt)
		break
	default:
		res += ": there are no additional information for node"
	}
	return res
}

func TreeToString() string {
	if Root != nil {
		return Root.TreeToString(1)
	}
	return "{ nil }"
}

// ROOT-RIGHT-LEFT scheme
func (node *Node) TreeToString(offset int) string {
	lowerOffsetStr := strings.Repeat("\t", offset-1)
	upperOffsetStr := strings.Repeat("\t", offset)
	nodeStr := upperOffsetStr + node.ToString() + "\n"
	var left string
	if node.Left != nil {
		left = node.Left.TreeToString(offset + 1)
	} else {
		left = upperOffsetStr + "{ nil }\n"
	}
	var right string
	if node.Right != nil {
		right = node.Right.TreeToString(offset + 1)
	} else {
		right = upperOffsetStr + "{ nil }\n"
	}
	return lowerOffsetStr + "{\n" + nodeStr + right + left + lowerOffsetStr + "}\n"
}
