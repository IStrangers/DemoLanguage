package vm

type BaseFunObject struct {
	BaseObject
	funDefinition string
	program       *Program
}

func (self *BaseFunObject) toLiteral() string {
	return self.funDefinition
}

type FunObject struct {
	BaseFunObject
}

func (self FunObject) vmCall(vm *VM, n int) {
	vm.pushCtx()
	vm.args = n
	vm.program = self.program
	vm.stack[vm.sp-1-n], vm.stack[vm.sp-2-n] = vm.stack[vm.sp-2-n], vm.stack[vm.sp-1-n]
	vm.pc = 0
}

type NativeFunCall struct {
	this Value
	args []Value
}

type NativeFunObject struct {
	BaseFunObject

	fun func(NativeFunCall) Value
}

func (self NativeFunObject) vmCall(vm *VM, n int) {
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
