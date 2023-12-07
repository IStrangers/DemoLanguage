package vm

type StackFrame struct {
	program      *Program
	functionName string
	pc           int
}

type StackFrameArray []StackFrame

type Exception struct {
	value Value
	stack StackFrameArray
}

type Runtime struct {
	globalObject *Object
}

func (self *Runtime) newObject() *Object {
	baseObject := &BaseObject{}
	baseObject.className = classObject
	baseObject.init()
	return &Object{baseObject}
}

func (self *Runtime) newFun(name string) *FunObject {
	funObject := &FunObject{}
	funObject.className = classFunction
	funObject.BaseObject.init()
	funObject.BaseObject.setProperty("name", StringValue(name))
	return funObject
}
