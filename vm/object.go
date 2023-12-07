package vm

const (
	classObject   = "Object"
	classFunction = "Function"
)

type ObjectImpl interface {
	getClassName() string
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
