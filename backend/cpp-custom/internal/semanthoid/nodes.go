package semanthoid

type Node interface {
	Create() error
}

type Procedure struct {
	NodeTypeLabel 		int
	Name 				string
	ParamsCount 		int
	ParamTypeLabels 	[]int
	Right 				*Node
	Left 				*Node
}

type Variable struct {
	NodeTypeLabel 	int
	Name 			string
	TypeLabel 		int
	Value 			int
	Left 			*Node
}



type Constant struct {
	NodeTypeLabel 	int
	Name 			string
	TypeLabel 		int
	Value 			int
	Left 			*Node
}

type For struct {
	NodeTypeLabel 	int
	Right 			*Node
	Left 			*Node
}