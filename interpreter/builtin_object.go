package interpreter

type BuiltinObject struct {
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
