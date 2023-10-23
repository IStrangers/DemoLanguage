package interpreter

import (
	"math"
	"strconv"
)

type ValueType int

func (self ValueType) String() string {
	switch self {
	case Skip:
		return "Skip"
	case Null:
		return "Null"
	case Boolean:
		return "Boolean"
	case Number:
		return "Number"
	case String:
		return "String"
	case Object:
		return "Object"
	case Reference:
		return "Reference"
	default:
		return ""
	}
}

const (
	_ ValueType = iota
	Skip
	Break
	Continue

	Null
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

var (
	Const_True_Value  = Value{Boolean, true}
	Const_False_Value = Value{Boolean, false}
)

func (self *Value) isSkip() bool {
	return self.valueType == Skip
}

func (self *Value) isBreak() bool {
	return self.value == Break
}

func (self *Value) isContinue() bool {
	return self.value == Continue
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

func (self *Value) isReferenceNumber() bool {
	if self.isReference() {
		reference := self.value.(ReferenceValue)
		value := reference.getValue()
		return value.isNumber()
	}
	return self.isNumber()
}

func (self *Value) getVal() any {
	if self.isReference() {
		reference := self.value.(ReferenceValue)
		return reference.getVal()
	}
	return self.value
}

func (self *Value) int64() int64 {
	switch v := self.getVal().(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	}
	panic("Unable to convert to int64")
}

func (self *Value) float64() float64 {
	switch v := self.getVal().(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	}
	panic("Unable to convert to float64")
}

func (self *Value) bool() bool {
	switch v := self.getVal().(type) {
	case bool:
		return v
	}
	panic("Unable to convert to bool")
}

func (self *Value) string() string {
	switch v := self.getVal().(type) {
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

func (self *Value) reference() ReferenceValue {
	if self.isReference() {
		return self.value.(ReferenceValue)
	}
	panic("Unable to convert to reference")
}

type ReferenceValue interface {
	getName() string
	getVal() any
	getValue() Value
	setValue(value Value)
}

type StashReference struct {
	name  string
	stash *Stash
}

func (self *StashReference) getName() string {
	return self.name
}

func (self *StashReference) getVal() any {
	value := self.getValue()
	return value.getVal()
}

func (self *StashReference) getValue() Value {
	value := self.stash.getValue(self.name)
	return value
}

func (self *StashReference) setValue(value Value) {
	self.stash.setValue(self.name, value)
}
