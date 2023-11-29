package vm

type BaseFunObject struct {
	BaseObject
	funDefinition string
	program       *Program
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
