package interpreter

type Arrayd struct {
	values []Value
}

func ArrayObject(values []Value) Value {
	object := Value{
		valueType: Object,
		value:     nil,
	}
	array := Arrayd{
		values: values,
	}
	object.value = Objectd{
		origin: array,
		propertys: map[string]Value{
			"get": {
				valueType: Function,
				value: Functiond{
					this: object,
					name: "get",
					callee: func(arguments ...Value) Value {
						return array.values[arguments[0].int64()]
					},
				},
			},
			"add": {
				valueType: Function,
				value: Functiond{
					this: object,
					name: "add",
					callee: func(arguments ...Value) Value {
						for _, argument := range arguments {
							array.values = append(array.values, argument)
						}
						return Const_Skip_Value
					},
				},
			},
			"remove": {
				valueType: Function,
				value: Functiond{
					this: object,
					name: "remove",
					callee: func(arguments ...Value) Value {
						index := arguments[0].int64()
						removeValue := array.values[index]
						array.values = append(array.values[:index], array.values[index+1:]...)
						return removeValue
					},
				},
			},
			"size": {
				valueType: Function,
				value: Functiond{
					this: object,
					name: "size",
					callee: func(arguments ...Value) Value {
						return NumberValue(len(array.values))
					},
				},
			},
		},
	}
	return object
}
