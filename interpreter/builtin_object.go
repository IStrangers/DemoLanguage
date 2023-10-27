package interpreter

import (
	"fmt"
	"strings"
)

type BuiltinObject struct {
}

func (self BuiltinObject) ofLiteral(object Objectd) string {
	var jsons []string
	for name, property := range object.propertys {
		jsons = append(jsons, fmt.Sprintf("%s: %s", name, property.ofLiteral()))
	}
	return fmt.Sprintf("{%s}", strings.Join(jsons, ","))
}

func (self BuiltinObject) getValue(object Objectd, property Value, args ...Value) Value {
	return object.getProperty(property.string())
}

func (self BuiltinObject) setValue(object Objectd, property Value, values ...Value) {
	object.setProperty(property.string(), values[0])
}

func BuiltinObjectObject(propertys map[string]Value) Value {
	object := Value{
		valueType: Object,
	}
	builtinArray := BuiltinObject{}
	object.value = Objectd{
		classObject: builtinArray,
		propertys:   propertys,
	}
	return object
}
