package vm

type ObjectImpl interface {
	getClassName() string
	getProperty(string) Value
	setProperty(string, Value)
	equals(objectImpl ObjectImpl) bool
}

type BaseObject struct {
	className    string
	valueMapping map[string]Value
}

func (self *BaseObject) getClassName() string {
	return self.className
}

func (self *BaseObject) getProperty(name string) Value {
	return self.valueMapping[name]
}

func (self *BaseObject) setProperty(name string, value Value) {
	self.valueMapping[name] = value
}

func (self *BaseObject) equals(objectImpl ObjectImpl) bool {
	return self == objectImpl
}
