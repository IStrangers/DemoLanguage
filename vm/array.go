package vm

import (
	"fmt"
	"strings"
)

type ArrayObject struct {
	BaseObject
	values ValueArray
	length uint32
}

func (self *ArrayObject) init() {
	self.BaseObject.valueMapping = arrayProps
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
		valueFormat := "%s"
		if value.isString() {
			valueFormat = "\"%s\""
		}
		literals = append(literals, fmt.Sprintf(valueFormat, value.toLiteral()))
	}
	return "[" + strings.Join(literals, ",") + "]"
}
