package vm

import "strconv"

type Value interface {
	isInt() bool
	isFloat() bool
	isString() bool

	toInt() int64
	toFloat() float64
	toString() string
}

type ValueArray []Value

type IntValue int64

func (self IntValue) isInt() bool {
	return true
}

func (self IntValue) isFloat() bool {
	return false
}

func (self IntValue) isString() bool {
	return false
}

func (self IntValue) toInt() int64 {
	return int64(self)
}

func (self IntValue) toFloat() float64 {
	return float64(self)
}

func (self IntValue) toString() string {
	return strconv.FormatInt(self.toInt(), 10)
}

func ToIntValue(value int64) IntValue {
	return IntValue(value)
}

type FloatValue float64

func (self FloatValue) isInt() bool {
	return false
}

func (self FloatValue) isFloat() bool {
	return true
}

func (self FloatValue) isString() bool {
	return false
}

func (self FloatValue) toInt() int64 {
	return int64(self)
}

func (self FloatValue) toFloat() float64 {
	return float64(self)
}

func (self FloatValue) toString() string {
	return strconv.FormatFloat(self.toFloat(), 'f', -1, 32)
}

func ToFloatValue(value float64) FloatValue {
	return FloatValue(value)
}

type StringValue string

func (self StringValue) isInt() bool {
	return false
}

func (self StringValue) isFloat() bool {
	return false
}

func (self StringValue) isString() bool {
	return true
}

func (self StringValue) toInt() int64 {
	v, _ := strconv.ParseInt(self.toString(), 0, 64)
	return v
}

func (self StringValue) toFloat() float64 {
	v, _ := strconv.ParseFloat(self.toString(), 64)
	return v
}

func (self StringValue) toString() string {
	return string(self)
}

func ToStringValue(value string) StringValue {
	return StringValue(value)
}
