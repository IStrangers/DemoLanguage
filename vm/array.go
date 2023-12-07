package vm

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
