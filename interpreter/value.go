package interpreter

import (
	"math"
	"strconv"
)

type ValueType int

const (
	_ ValueType = iota
	Skip
	NULL
	Boolean
	Number
	String
	Object
	Reference
)

type Value struct {
	valueType ValueType
	value     any
}

func (self *Value) isSkip() bool {
	return self.valueType == Skip
}

func (self *Value) isBoolean() bool {
	return self.valueType == Boolean
}

func (self *Value) isNumber() bool {
	return self.valueType == Number
}

func (self *Value) isString() bool {
	return self.valueType == String
}

func (self *Value) isObject() bool {
	return self.valueType == Object
}

func (self *Value) isReference() bool {
	return self.valueType == Reference
}

func (self *Value) getValue() any {
	if self.isReference() {
		reference := self.value.(Value)
		return reference.getValue()
	}
	return self.value
}

func (self *Value) int64() int64 {
	switch v := self.value.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}

func (self *Value) float64() float64 {
	switch v := self.value.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}

func (self *Value) bool() bool {
	switch v := self.value.(type) {
	case bool:
		return v
	default:
		return false
	}
}

func (self *Value) string() string {
	switch v := self.value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return floatToString(v, 32)
	case string:
		return v
	default:
		return ""
	}
}

func floatToString(value float64, bitSize int) string {
	if math.IsNaN(value) {
		return "NaN"
	} else if math.IsInf(value, 0) {
		if math.Signbit(value) {
			return "-Infinity"
		}
		return "Infinity"
	}
	return strconv.FormatFloat(value, 'f', -1, bitSize)
}
