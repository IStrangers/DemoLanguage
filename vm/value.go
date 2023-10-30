package vm

import "strconv"

type Value interface {
	isInt() bool
	isFloat() bool
	isString() bool
	isBool() bool
	isNull() bool

	toInt() int64
	toFloat() float64
	toString() string
	toBool() bool
}

type ValueArray []Value

var (
	Const_Bool_True_Value  = BoolValue(true)
	Const_Bool_False_Value = BoolValue(false)
	Const_Null_Value       = NullValue{}
)

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

func (self IntValue) isBool() bool {
	return false
}

func (self IntValue) isNull() bool {
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

func (self IntValue) toBool() bool {
	if self.toInt() > 0 {
		return true
	}
	return false
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

func (self FloatValue) isBool() bool {
	return false
}

func (self FloatValue) isNull() bool {
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

func (self FloatValue) toBool() bool {
	if self.toFloat() > 0.0 {
		return true
	}
	return false
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

func (self StringValue) isBool() bool {
	return false
}

func (self StringValue) isNull() bool {
	return false
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

func (self StringValue) toBool() bool {
	if len(self.toString()) > 0 {
		return true
	}
	return false
}

func ToStringValue(value string) StringValue {
	return StringValue(value)
}

type BoolValue bool

func (self BoolValue) isInt() bool {
	return false
}

func (self BoolValue) isFloat() bool {
	return false
}

func (self BoolValue) isString() bool {
	return false
}

func (self BoolValue) isBool() bool {
	return true
}

func (self BoolValue) isNull() bool {
	return false
}

func (self BoolValue) toInt() int64 {
	if self.toBool() {
		return 1
	}
	return 0
}

func (self BoolValue) toFloat() float64 {
	if self.toBool() {
		return 1.0
	}
	return 0.0
}

func (self BoolValue) toString() string {
	return strconv.FormatBool(self.toBool())
}

func (self BoolValue) toBool() bool {
	return bool(self)
}

type NullValue struct{}

func (self NullValue) isInt() bool {
	return false
}

func (self NullValue) isFloat() bool {
	return false
}

func (self NullValue) isString() bool {
	return false
}

func (self NullValue) isBool() bool {
	return false
}

func (self NullValue) isNull() bool {
	return true
}

func (self NullValue) toInt() int64 {
	return 0
}

func (self NullValue) toFloat() float64 {
	return 0.0
}

func (self NullValue) toString() string {
	return "null"
}

func (self NullValue) toBool() bool {
	return false
}
