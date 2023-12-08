package vm

import (
	"fmt"
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
	globalObject *Object
}

func CreateRuntime() *Runtime {
	runtime := &Runtime{
		globalObject: &Object{self: &BaseObject{
			valueMapping: map[string]Value{
				"println": Object{&NativeFunObject{fun: func(call NativeFunCall) Value {
					var literals []any
					for _, arg := range call.args {
						literals = append(literals, arg.toLiteral())
					}
					fmt.Fprintln(os.Stdout, literals...)
					return nil
				}}},
			},
		}},
	}
	return runtime
}

func (self *Runtime) newObject() *Object {
	baseObject := &BaseObject{}
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

func (self *Runtime) newFun(name string) *FunObject {
	funObject := &FunObject{}
	funObject.className = classFunction
	funObject.init()
	funObject.BaseObject.setProperty("name", StringValue(name))
	return funObject
}
