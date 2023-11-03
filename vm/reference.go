package vm

type Ref interface {
	name() string
	get() Value
	set(Value)
}

type ObjectRef struct {
	refObject *Object
	refName   string
}

func (self *ObjectRef) name() string {
	return self.refName
}

func (self *ObjectRef) get() Value {
	return self.refObject.self.getProperty(self.refName)
}

func (self *ObjectRef) set(value Value) {
	self.refObject.self.setProperty(self.refName, value)
}
