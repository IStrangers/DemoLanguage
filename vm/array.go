package vm

import "strings"

type ArrayObject struct {
	BaseObject
	values ValueArray
	length uint32
}

func (self *ArrayObject) init() {
	self.BaseObject.init()
}

func (self *ArrayObject) getValueByIndex(prop IntValue, defaultValue Value) Value {
	index := uint32(prop.toInt())
	if index < 0 || index >= self.length {
		return defaultValue
	}
	return self.values[index]
}

func (self *ArrayObject) toLiteral() string {
	var literals []string
	for _, value := range self.values {
		literals = append(literals, value.toLiteral())
	}
	return "[" + strings.Join(literals, ",") + "]"
}
