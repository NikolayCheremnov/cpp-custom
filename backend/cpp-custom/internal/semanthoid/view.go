package semanthoid

import (
	"strconv"
	"strings"
)

// STACK VIEWS
func (s *ProcStack) ToString() string {
	res := ""
	for i := 0; i < len(*s); i++ {
		res += (*s)[i].ToString() + " "
	}
	return res
}

func (c *ProcContext) ToString() string {
	return "[" + c.Proc.Identifier + "]"
}

// PROCEDURES LIST VIEWS

func (proc *ProcNode) ToString() string {
	res := "void " + proc.Identifier + "("
	for paramIndex, paramIdentifier := range proc.ParamsIdentifiers {
		paramTypeStr := "int"
		if proc.ParamsTypesLabels[paramIndex] == BoolType {
			paramTypeStr = "bool"
		}
		res += paramTypeStr + " " + paramIdentifier + ", "
	}
	if proc.ParamsCount > 0 {
		res = res[:len(res)-2] + ")"
	} else {
		res += ")"
	}
	return res
}

func ProcListToString() string {
	if ProcRoot == nil {
		return "nil"
	}
	res := ""
	for next := ProcRoot; next != nil; next = next.Next {
		res += next.ToString() + "\n"
	}
	return res
}

// TREE VIEWS

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
