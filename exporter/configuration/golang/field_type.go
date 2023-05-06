package golang

import "github.com/kercylan98/minotaur/exporter/configuration"

func GetFieldType(fieldType string) configuration.FieldType {
	switch fieldType {
	case "string":
		return configuration.FieldTypeString
	case "int":
		return configuration.FieldTypeInt
	case "int8":
		return configuration.FieldTypeInt8
	case "int16":
		return configuration.FieldTypeInt16
	case "int32":
		return configuration.FieldTypeInt32
	case "int64":
		return configuration.FieldTypeInt64
	case "uint":
		return configuration.FieldTypeUint
	case "uint8":
		return configuration.FieldTypeUint8
	case "uint16":
		return configuration.FieldTypeUint16
	case "uint32":
		return configuration.FieldTypeUint32
	case "uint64":
		return configuration.FieldTypeUint64
	case "float32":
		return configuration.FieldTypeFloat32
	case "float64":
		return configuration.FieldTypeFloat64
	case "byte":
		return configuration.FieldTypeByte
	case "bool":
		return configuration.FieldTypeBool
	case "rune":
		return configuration.FieldTypeRune
	case "[]string":
		return configuration.FieldTypeSliceString
	case "[]int":
		return configuration.FieldTypeSliceInt
	case "[]int8":
		return configuration.FieldTypeSliceInt8
	case "[]int16":
		return configuration.FieldTypeSliceInt16
	case "[]int32":
		return configuration.FieldTypeSliceInt32
	case "[]int64":
		return configuration.FieldTypeSliceInt64
	case "[]uint":
		return configuration.FieldTypeSliceUint
	case "[]uint8":
		return configuration.FieldTypeSliceUint8
	case "[]uint16":
		return configuration.FieldTypeSliceUint16
	case "[]uint32":
		return configuration.FieldTypeSliceUint32
	case "[]uint64":
		return configuration.FieldTypeSliceUint64
	case "[]float32":
		return configuration.FieldTypeSliceFloat32
	case "[]float64":
		return configuration.FieldTypeSliceFloat64
	case "[]byte":
		return configuration.FieldTypeSliceByte
	case "[]bool":
		return configuration.FieldTypeSliceBool
	case "[]rune":
		return configuration.FieldTypeSliceRune
	default:
		return configuration.FieldTypeString
	}
}

func GetFieldTypeName(fieldType configuration.FieldType) string {
	switch fieldType {
	case configuration.FieldTypeString:
		return "string"
	case configuration.FieldTypeInt:
		return "int"
	case configuration.FieldTypeInt8:
		return "int8"
	case configuration.FieldTypeInt16:
		return "int16"
	case configuration.FieldTypeInt32:
		return "int32"
	case configuration.FieldTypeInt64:
		return "int64"
	case configuration.FieldTypeUint:
		return "uint"
	case configuration.FieldTypeUint8:
		return "uint8"
	case configuration.FieldTypeUint16:
		return "uint16"
	case configuration.FieldTypeUint32:
		return "uint32"
	case configuration.FieldTypeUint64:
		return "uint64"
	case configuration.FieldTypeFloat32:
		return "float32"
	case configuration.FieldTypeFloat64:
		return "float64"
	case configuration.FieldTypeByte:
		return "byte"
	case configuration.FieldTypeBool:
		return "bool"
	case configuration.FieldTypeRune:
		return "rune"
	case configuration.FieldTypeSliceString:
		return "[]string"
	case configuration.FieldTypeSliceInt:
		return "[]int"
	case configuration.FieldTypeSliceInt8:
		return "[]int8"
	case configuration.FieldTypeSliceInt16:
		return "[]int16"
	case configuration.FieldTypeSliceInt32:
		return "[]int32"
	case configuration.FieldTypeSliceInt64:
		return "[]int64"
	case configuration.FieldTypeSliceUint:
		return "[]uint"
	case configuration.FieldTypeSliceUint8:
		return "[]uint8"
	case configuration.FieldTypeSliceUint16:
		return "[]uint16"
	case configuration.FieldTypeSliceUint32:
		return "[]uint32"
	case configuration.FieldTypeSliceUint64:
		return "[]uint64"
	case configuration.FieldTypeSliceFloat32:
		return "[]float32"
	case configuration.FieldTypeSliceFloat64:
		return "[]float64"
	case configuration.FieldTypeSliceByte:
		return "[]byte"
	case configuration.FieldTypeSliceBool:
		return "[]bool"
	case configuration.FieldTypeSliceRune:
		return "[]rune"
	default:
		return ""
	}
}
