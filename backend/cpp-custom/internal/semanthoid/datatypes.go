package semanthoid

type DataTypeValue struct {
	DataAsInt  int
	DataAsBool int
}

func GetDefaultDataValue() *DataTypeValue {
	return &DataTypeValue{0, 0}
}

func ConvertToDataTypeValue(dataTypeLabel int, value int) *DataTypeValue {
	var dataValue *DataTypeValue
	switch dataTypeLabel {
	case IntType:
		dataValue = &DataTypeValue{DataAsInt: value, DataAsBool: IntToBool(value)}
		break
	case BoolType:
		dataValue = &DataTypeValue{DataAsInt: IntToBool(value), DataAsBool: IntToBool(value)}
		break
	}
	return dataValue
}

func IntToBool(dataAsInt int) int {
	if dataAsInt != 0 {
		return 1
	}
	return 0
}

func GoBoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
