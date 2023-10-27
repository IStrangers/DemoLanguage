package interpreter

import "strings"

const (
	BuiltinArray_Get_Method    = "get"
	BuiltinArray_Add_Method    = "add"
	BuiltinArray_Remove_Method = "remove"
	BuiltinArray_Size_Method   = "size"
)

type BuiltinArray struct {
	values []Value
}

func (self BuiltinArray) ofLiteral(object Objectd) string {
	var jsons []string
	for _, value := range self.values {
		jsons = append(jsons, value.ofLiteral())
	}
	return "[" + strings.Join(jsons, ",") + "]"
}

func (self BuiltinArray) getValue(object Objectd, property Value, args ...Value) Value {
	if property.isString() {
		return object.getProperty(property.string())
	}
	getMethod := object.getProperty(BuiltinArray_Get_Method)
	return getMethod.functiond().call(property)
}

func (self BuiltinArray) setValue(object Objectd, property Value, values ...Value) {
	addMethod := object.getProperty(BuiltinArray_Add_Method)
	addMethod.functiond().call(values...)
}

func BuiltinArrayObject(values []Value) Value {
	object := Value{
		valueType: Object,
	}
	builtinArray := BuiltinArray{
		values: values,
	}
	object.value = Objectd{
		classObject: builtinArray,
		propertys: map[string]Value{
			BuiltinArray_Get_Method: FunctionValue(Functiond{
				this: object,
				name: BuiltinArray_Get_Method,
				callee: func(arguments ...Value) Value {
					return builtinArray.values[arguments[0].int64()]
				},
			}),
			BuiltinArray_Add_Method: FunctionValue(Functiond{
				this: object,
				name: BuiltinArray_Add_Method,
				callee: func(arguments ...Value) Value {
					for _, argument := range arguments {
						builtinArray.values = append(builtinArray.values, argument)
					}
					return Const_Skip_Value
				},
			}),
			BuiltinArray_Remove_Method: FunctionValue(Functiond{
				this: object,
				name: BuiltinArray_Remove_Method,
				callee: func(arguments ...Value) Value {
					index := arguments[0].int64()
					removeValue := builtinArray.values[index]
					builtinArray.values = append(builtinArray.values[:index], builtinArray.values[index+1:]...)
					return removeValue
				},
			}),
			BuiltinArray_Size_Method: FunctionValue(Functiond{
				this: object,
				name: BuiltinArray_Size_Method,
				callee: func(arguments ...Value) Value {
					return NumberValue(len(builtinArray.values))
				},
			}),
		},
	}
	return object
}
