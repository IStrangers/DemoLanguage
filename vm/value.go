package vm

import "strconv"

type Value interface {
	isInt() bool
	isFloat() bool
	isString() bool
	isBool() bool
	isNull() bool
	isObject() bool

	toInt() int64
	toFloat() float64
	toString() string
	toBool() bool
	toObject() *Object

	equals(Value) bool
	sameAs(Value) bool
}

type ValueArray []Value

func (self ValueArray) findIndex(value Value) int {
	for index, v := range self {
		if v.sameAs(value) {
			return index
		}
	}
	return -1
}

func (self *ValueArray) add(values ...Value) {
	*self = append(*self, values...)
}

func (self *ValueArray) size() int {
	return len(*self)
}

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

func (self IntValue) isObject() bool {
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
	return self.toInt() > 0
}

func (self IntValue) toObject() *Object {
	return nil
}

func (self IntValue) equals(value Value) bool {
	if self.sameAs(value) {
		return true
	}
	return self.toFloat() == value.toFloat()
}

func (self IntValue) sameAs(value Value) bool {
	if value.isInt() {
		return self == value
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

func (self FloatValue) isObject() bool {
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
	return self.toFloat() > 0.0
}

func (self FloatValue) toObject() *Object {
	return nil
}

func (self FloatValue) equals(value Value) bool {
	if self.sameAs(value) {
		return true
	}
	return self.toFloat() == value.toFloat()
}

func (self FloatValue) sameAs(value Value) bool {
	if value.isFloat() {
		return self == value
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

func (self StringValue) isObject() bool {
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
	return len(self.toString()) > 0
}

func (self StringValue) toObject() *Object {
	return nil
}

func (self StringValue) equals(value Value) bool {
	if self.sameAs(value) {
		return true
	}
	return self.toString() == value.toString()
}

func (self StringValue) sameAs(value Value) bool {
	if value.isString() {
		return self == value
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

func (self BoolValue) isObject() bool {
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

func (self BoolValue) toObject() *Object {
	return nil
}

func (self BoolValue) equals(value Value) bool {
	if self.sameAs(value) {
		return true
	}
	return self.toBool() == value.toBool()
}

func (self BoolValue) sameAs(value Value) bool {
	if value.isBool() {
		return self == value
	}
	return false
}

func ToBooleanValue(value bool) BoolValue {
	if value {
		return Const_Bool_True_Value
	}
	return Const_Bool_False_Value
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

func (self NullValue) isObject() bool {
	return false
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

func (self NullValue) toObject() *Object {
	return nil
}

func (self NullValue) equals(value Value) bool {
	return self.sameAs(value)
}

func (self NullValue) sameAs(value Value) bool {
	return value.isNull()
}

type Object struct {
	self ObjectImpl
}

func (self Object) isInt() bool {
	return false
}

func (self Object) isFloat() bool {
	return false
}

func (self Object) isString() bool {
	return false
}

func (self Object) isBool() bool {
	return false
}

func (self Object) isNull() bool {
	return false
}

func (self Object) isObject() bool {
	return true
}

func (self Object) toInt() int64 {
	return 0
}

func (self Object) toFloat() float64 {
	return 0.0
}

func (self Object) toString() string {
	return "Object"
}

func (self Object) toBool() bool {
	return self.self != nil
}

func (self Object) toObject() *Object {
	return &self
}

func (self Object) equals(value Value) bool {
	return self.sameAs(value)
}

func (self Object) sameAs(value Value) bool {
	if value.isObject() {
		return self == value || self.self.equals(value.(Object).self)
	}
	return false
}

func (self *Object) getOrDefault(prop Value, defaultValue Value) Value {
	if prop.isInt() {
		return self.self.getValueByIndex(prop.(IntValue), defaultValue)
	}
	return self.self.getPropertyOrDefault(prop.toString(), defaultValue)
}
