package vm

import (
	"fmt"
	"strings"
)

const (
	thisBindingName = "this"

	classObject   = "Object"
	classGlobal   = "Global"
	classArray    = "Array"
	classFunction = "Function"
)

type ObjectImpl interface {
	getClassName() string
	toLiteral() string
	getValueByIndex(IntValue, Value) Value
	getProperty(string) Value
	getPropertyOrDefault(string, Value) Value
	setProperty(string, Value)
	equals(objectImpl ObjectImpl) bool
	vmCall(vm *VM, n int)
}

type BaseObject struct {
	className    string
	valueMapping map[string]Value
}

func (self *BaseObject) init() {
	self.valueMapping = make(map[string]Value)
}

func (self *BaseObject) getClassName() string {
	return self.className
}

func (self *BaseObject) toLiteral() string {
	className := self.getClassName()
	if className != classObject {
		return fmt.Sprintf("[Object %s]", className)
	}
	var literals []string
	for name, value := range self.valueMapping {
		valueFormat := "%s"
		if value.isString() {
			valueFormat = "\"%s\""
		}
		literals = append(literals, fmt.Sprintf("%s: %s", name, fmt.Sprintf(valueFormat, value.toLiteral())))
	}
	return fmt.Sprintf("{%s}", strings.Join(literals, ","))
}

func (self *BaseObject) getValueByIndex(prop IntValue, defaultValue Value) Value {
	return self.getPropertyOrDefault(prop.toString(), defaultValue)
}

func (self *BaseObject) getProperty(name string) Value {
	return self.valueMapping[name]
}

func (self *BaseObject) getPropertyOrDefault(name string, defaultValue Value) Value {
	value := self.getProperty(name)
	if value == nil {
		return defaultValue
	}
	return value
}

func (self *BaseObject) setProperty(name string, value Value) {
	self.valueMapping[name] = value
}

func (self *BaseObject) equals(objectImpl ObjectImpl) bool {
	return self == objectImpl
}

func (self *BaseObject) vmCall(vm *VM, n int) {
	//wait adjust
	panic("Not a function: " + self.className)
}

type ClassObject struct {
	BaseObject
	classDefinition string
	constructors    []*ClassFunObject
}

func (self *ClassObject) findConstructor(argNum int) *ClassFunObject {
	for _, constructor := range self.constructors {
		if constructor.argNum <= argNum {
			return constructor
		}
	}
	return nil
}
