package interpreter

import (
	"fmt"
	"os"
)

const (
	BuiltinGlobal_Println_Method = "println"
)

type BuiltinGlobal struct {
	values []Value
}

func (self BuiltinGlobal) json(object Objectd) string {
	return "@Class BuiltinGlobal"
}

func (self BuiltinGlobal) getValue(object Objectd, property Value, args ...Value) Value {
	return object.getProperty(property.string())
}

func (self BuiltinGlobal) setValue(object Objectd, property Value, values ...Value) {
	object.setProperty(property.string(), values[0])
}

func BuiltinGlobalObject() Objectd {
	builtinGlobal := BuiltinGlobal{}
	object := Objectd{
		classObject: builtinGlobal,
		propertys: map[string]Value{
			BuiltinGlobal_Println_Method: FunctionValue(Functiond{
				name: BuiltinGlobal_Println_Method,
				callee: func(arguments ...Value) Value {
					fmt.Fprintln(os.Stdout, toVals(arguments)...)
					return Const_Skip_Value
				},
			}),
		},
	}
	return object
}

func toVals(values []Value) []any {
	var vals []any
	for _, value := range values {
		vals = append(vals, value.json())
	}
	return vals
}
