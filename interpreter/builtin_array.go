package interpreter

const (
	Get_Method    = "get"
	Add_Method    = "add"
	Remove_Method = "remove"
	Size_Method   = "size"
)

type BuiltinArray struct {
	values []Value
}

func (self BuiltinArray) getValue(object Objectd, property Value, args ...Value) Value {
	getMethod := object.getProperty(Get_Method)
	return getMethod.functiond().call(property)
}

func (self BuiltinArray) setValue(object Objectd, property Value, values ...Value) {
	addMethod := object.getProperty(Add_Method)
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
			Get_Method: {
				valueType: Function,
				value: Functiond{
					this: object,
					name: Get_Method,
					callee: func(arguments ...Value) Value {
						return builtinArray.values[arguments[0].int64()]
					},
				},
			},
			Add_Method: {
				valueType: Function,
				value: Functiond{
					this: object,
					name: Add_Method,
					callee: func(arguments ...Value) Value {
						for _, argument := range arguments {
							builtinArray.values = append(builtinArray.values, argument)
						}
						return Const_Skip_Value
					},
				},
			},
			Remove_Method: {
				valueType: Function,
				value: Functiond{
					this: object,
					name: Remove_Method,
					callee: func(arguments ...Value) Value {
						index := arguments[0].int64()
						removeValue := builtinArray.values[index]
						builtinArray.values = append(builtinArray.values[:index], builtinArray.values[index+1:]...)
						return removeValue
					},
				},
			},
			Size_Method: {
				valueType: Function,
				value: Functiond{
					this: object,
					name: Size_Method,
					callee: func(arguments ...Value) Value {
						return NumberValue(len(builtinArray.values))
					},
				},
			},
		},
	}
	return object
}
