package interpreter

import (
	"fmt"
	"strings"
)

type BuiltinObject struct {
}

func (self BuiltinObject) ofLiteral(object Objectd) string {
	var literals []string
	for name, property := range object.propertys {
		valueFormat := "%s"
		if property.isString() {
			valueFormat = "\"%s\""
		}
		literals = append(literals, fmt.Sprintf("%s: %s", name, fmt.Sprintf(valueFormat, property.ofLiteral())))
	}
	return fmt.Sprintf("{%s}", strings.Join(literals, ","))
}

func (self BuiltinObject) getValue(object Objectd, property Value, args ...Value) Value {
	return object.getProperty(property.string())
}

func (self BuiltinObject) setValue(object Objectd, property Value, values ...Value) {
	object.setProperty(property.string(), values[0])
}

func BuiltinObjectObject(propertys map[string]Value) Value {
	builtinArray := BuiltinObject{}
	object := Objectd{
		classObject: builtinArray,
		propertys:   propertys,
	}
	return ObjectValue(object)
}
