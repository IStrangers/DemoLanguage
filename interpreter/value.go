package interpreter

import (
	"math"
	"strconv"
)

type ValueType int

const (
	_ ValueType = iota
	Skip
	Break
	Continue
	Return

	Null
	Boolean
	Number
	String
	Object
	Function
	MultipleValue
	Reference
)

type Value struct {
	valueType ValueType
	value     any
}

var (
	Const_Skip_Value     = Value{Skip, nil}
	Const_Break_Value    = Value{Break, nil}
	Const_Continue_Value = Value{Continue, nil}
	Const_Return_Value   = Value{Return, nil}
	Const_Null_Value     = Value{Null, nil}
	Const_True_Value     = Value{Boolean, true}
	Const_False_Value    = Value{Boolean, false}
)

func (self *Value) isResult() bool {
	return !self.isSkip() && !self.isBreak() && !self.isContinue()
}

func (self *Value) isSkip() bool {
	return self.valueType == Skip
}

func (self *Value) isBreak() bool {
	return self.valueType == Break
}

func (self *Value) isContinue() bool {
	return self.valueType == Continue
}

func (self *Value) isReturn() bool {
	return self.valueType == Return
}

func (self *Value) isBoolean() bool {
	flatValue := self.flatResolve()
	return flatValue.valueType == Boolean
}

func (self *Value) isNumber() bool {
	flatValue := self.flatResolve()
	return flatValue.valueType == Number
}

func (self *Value) isString() bool {
	flatValue := self.flatResolve()
	return flatValue.valueType == String
}

func (self *Value) isObject() bool {
	flatValue := self.flatResolve()
	return flatValue.valueType == Object
}

func (self *Value) isFunction() bool {
	flatValue := self.flatResolve()
	return flatValue.valueType == Function
}

func (self *Value) isMultipleValue() bool {
	flatValue := self.flatResolve()
	return flatValue.valueType == MultipleValue
}

func (self *Value) isReferenced() bool {
	return self.valueType == Reference
}

func (self Value) flatResolve() Value {
	if self.isReferenced() {
		reference := self.referenced()
		value := reference.getValue()
		if value.isReferenced() {
			value = value.flatResolve()
		}
		return value
	}
	return self
}

func (self *Value) getVal() any {
	flatValue := self.flatResolve()
	if flatValue.isReturn() {
		ofValue := flatValue.ofValue()
		return ofValue.getVal()
	}
	return flatValue.value
}

func (self *Value) ofValue() Value {
	switch v := self.value.(type) {
	case Value:
		return v
	}
	panic("Unable to convert to value")
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
	}
	panic("Unable to convert to string")
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

func (self *Value) objectd() Objectd {
	if self.isObject() {
		return self.value.(Objectd)
	}
	panic("Unable to convert to object")
}

func (self *Value) functiond() Functiond {
	if self.isFunction() {
		return self.value.(Functiond)
	}
	panic("Unable to convert to function")
}

func (self *Value) referenced() Referenced {
	if self.isReferenced() {
		return self.value.(Referenced)
	}
	panic("Unable to convert to reference")
}

type Objectd struct {
	propertys map[string]Value
}

func (self *Objectd) getProperty(name string) Value {
	value, exists := self.propertys[name]
	if !exists {
		return Const_Null_Value
	}
	return value
}

func (self *Objectd) setProperty(name string, value Value) {
	self.propertys[name] = value
}

func (self *Objectd) contains(name string) bool {
	_, exists := self.propertys[name]
	return exists
}

type Functiond interface {
	getName() string
	call(arguments ...Value) Value
}

type GlobalFunctiond struct {
	name   string
	callee func(arguments ...Value) Value
}

func (self *GlobalFunctiond) getName() string {
	return self.name
}

func (self *GlobalFunctiond) call(arguments ...Value) Value {
	return self.callee(arguments...)
}

type Referenced interface {
	getName() string
	getVal() any
	getValue() Value
	setValue(value Value)
}

type StashReferenced struct {
	name  string
	stash *Stash
}

func (self *StashReferenced) getName() string {
	return self.name
}

func (self *StashReferenced) getVal() any {
	value := self.getValue()
	return value.getVal()
}

func (self *StashReferenced) getValue() Value {
	value := self.stash.getValue(self.name)
	return value
}

func (self *StashReferenced) setValue(value Value) {
	self.stash.setValue(self.name, value)
}
