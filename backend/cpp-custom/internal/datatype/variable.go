package datatype

type Variable struct {
	IsMutable bool
	Type      Type
}

func (v *Variable) VariableToString() string {
	if !v.IsMutable {
		return "const " + v.Type.FullName
	}
	return v.Type.FullName
}
