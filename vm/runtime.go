package vm

import (
	"fmt"
	"io"
	"os"
)

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
	global       *Global
	globalObject *Object
	vm           *VM
}

var RuntimePrintWrite io.Writer = os.Stdout

func CreateRuntime() *Runtime {
	runtime := &Runtime{
		global: &Global{},
		globalObject: &Object{self: &BaseObject{
			className: classGlobal,
			valueMapping: map[string]Value{
				"println": Object{&NativeFunObject{fun: func(call NativeFunCall) Value {
					var literals []any
					for _, arg := range call.args {
						literals = append(literals, arg.toLiteral())
					}
					fmt.Fprintln(RuntimePrintWrite, literals...)
					return nil
				}}},
			},
		}},
	}
	runtime.vm = &VM{
		runtime:          runtime,
		sb:               -1,
		maxCallStackSize: 999,
	}
	return runtime
}

func (self *Runtime) newObject() *Object {
	baseObject := &BaseObject{}
	baseObject.objectType = normalObject
	baseObject.className = classObject
	baseObject.init()
	return &Object{self: baseObject}
}

func (self *Runtime) newObjectByClass(className string) *Object {
	baseObject := &BaseObject{}
	baseObject.objectType = normalObject
	baseObject.className = className
	baseObject.init()
	return &Object{self: baseObject}
}

func (self *Runtime) newClassObject() *Object {
	baseObject := &ClassObject{}
	baseObject.objectType = classDefinitionObject
	baseObject.className = classObject
	baseObject.init()
	return &Object{baseObject}
}

func (self *Runtime) newArray(values ValueArray) *Object {
	arrayObject := &ArrayObject{}
	arrayObject.className = classArray
	arrayObject.values = values
	arrayObject.length = uint32(values.size())
	arrayObject.init()
	return &Object{arrayObject}
}

func (self *Runtime) newFun(name string, length int) *FunObject {
	funObject := &FunObject{}
	funObject.className = classFunction
	funObject.init()
	funObject.setProperty("name", ToStringValue(name))
	funObject.setProperty("length", ToIntValue(int64(length)))
	return funObject
}

func (self *Runtime) newClassFun(name string, length int) *ClassFunObject {
	funObject := &ClassFunObject{}
	funObject.className = classFunction
	funObject.init()
	funObject.setProperty("name", ToStringValue(name))
	funObject.setProperty("length", ToIntValue(int64(length)))
	return funObject
}

func (self *Runtime) createReferenceError(msg string) *Object {
	if self.global.referenceError == nil {
		self.global.referenceError = self.newObject()
	}
	self.global.referenceError.self.setProperty("message", ToStringValue(msg))
	return self.global.referenceError
}

func (self *Runtime) newReferenceError(name string) Value {
	return self.createReferenceError(fmt.Sprintf("ReferenceError: '%s' is not defined", name))
}
