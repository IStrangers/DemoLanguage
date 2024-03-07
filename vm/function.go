package vm

type BaseFunObject struct {
	BaseObject
	funDefinition string
	argNum        int
	program       *Program
	stash         *Stash
}

func (self *BaseFunObject) toLiteral() string {
	return self.funDefinition
}

type FunObject struct {
	BaseFunObject
}

func (self *FunObject) vmCall(vm *VM, n int) {
	vm.pushCtx()
	vm.args = n
	vm.program = self.program
	vm.stash = self.stash
	vm.stack[vm.sp-1-n], vm.stack[vm.sp-2-n] = vm.stack[vm.sp-2-n], vm.stack[vm.sp-1-n]
	vm.pc = 0
}

type ClassFunObject struct {
	BaseFunObject
	initProgram *Program
}

func (self *ClassFunObject) vmCall(vm *VM, n int) {
	//wait adjust
	panic("Class constructor cannot be invoked without 'new'")
}

func (self *ClassFunObject) classConstruct(runtime *Runtime, args []Value) *Object {
	thisObj := runtime.newClassObject()
	return self.construct(runtime, thisObj, args)
}

func (self *ClassFunObject) instanceConstruct(runtime *Runtime, args []Value) *Object {
	thisObj := runtime.newObjectByClass(self.getProperty("name").toString())
	return self.construct(runtime, thisObj, args)
}

func (self *ClassFunObject) construct(runtime *Runtime, thisObj *Object, args []Value) *Object {
	self.initObject(runtime, thisObj)
	self.call(runtime, thisObj, args)
	return thisObj
}

func (self *ClassFunObject) initObject(runtime *Runtime, thisObj *Object) {
	if self.initProgram == nil {
		return
	}
	vm := runtime.vm
	vm.pushCtx()
	vm.program = self.initProgram
	vm.stash = self.stash

	vm.sb = vm.sp
	vm.push(thisObj)
	vm.pc = 0
	ex := vm.runTry()
	vm.popCtx()
	if ex != nil {
		panic(ex)
	}
	vm.sp -= 2
}

func (self *ClassFunObject) call(runtime *Runtime, thisObj *Object, args []Value) (Value, *Exception) {
	if self.program == nil {
		return nil, nil
	}
	vm := runtime.vm
	vm.stack.expand(vm.sp + len(args) + 1)
	vm.sp++
	vm.stack[vm.sp] = thisObj
	vm.sp++
	for _, arg := range args {
		if arg != nil {
			vm.stack[vm.sp] = arg
		} else {
			vm.stack[vm.sp] = Const_Null_Value
		}
		vm.sp++
	}

	vm.pushTryFrame(-2, -1)
	defer vm.popTryFrame()

	var needPop bool
	if vm.program != nil {
		vm.pushCtx()
		vm.callStack = append(vm.callStack, Context{pc: -2})
		needPop = true
	} else {
		vm.pc = -2
		vm.pushCtx()
	}

	vm.args = len(args)
	vm.program = self.program
	vm.stash = self.stash
	vm.pc = 0
	for {
		ex := vm.runTryInner()
		if ex != nil {
			return nil, ex
		}
		if vm.halted() {
			break
		}
	}
	if needPop {
		vm.popCtx()
	}

	return vm.pop(), nil
}

type NativeFunCall struct {
	this Value
	args []Value
}

type NativeFunObject struct {
	BaseFunObject

	fun func(NativeFunCall) Value
}

func (self *NativeFunObject) vmCall(vm *VM, n int) {
	vm.pushCtx()
	vm.program = nil
	vm.sb = vm.sp - n
	value := self.fun(NativeFunCall{
		this: vm.stack[vm.sp-n-2],
		args: vm.stack[vm.sp-n : vm.sp],
	})
	if value == nil {
		value = Const_Null_Value
	}
	vm.stack[vm.sp-n-2] = value
	vm.popCtx()
	vm.sp -= n + 1
	vm.pc++
}
