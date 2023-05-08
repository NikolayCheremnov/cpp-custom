package datatype

import "strconv"

type Value struct {
	DataAsInt  int64
	DataAsBool bool
}

func NewDefaultDataValue() Value {
	return Value{0, false}
}

func NewIntDataValue(value int64) Value {
	return Value{value, value != 0}
}

func NewBoolDataValue(value bool) Value {
	var intValue int64 = 0
	if value {
		intValue = 1
	}
	return Value{intValue, value}
}

func (v *Value) ValueAsString() string {
	return strconv.FormatInt(v.DataAsInt, 10) + "(" + strconv.FormatBool(v.DataAsBool) + ")"
}
