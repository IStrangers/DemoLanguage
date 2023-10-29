package vm

type Instruction interface {
	exec(*VM)
}

type InstructionArray []Instruction

var (
	Addition _Addition
	Subtract _Subtract
	Multiply _Multiply
	Divide   _Divide
)

type LoadValue int

func (self LoadValue) exec(vm *VM) {
	vm.push(vm.getValue(int(self)))
	vm.pc++
}

type _Addition struct{}

func (self _Addition) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isString() || right.isString() {
		value = ToStringValue(left.toString() + right.toString())
	} else if left.isFloat() || right.isFloat() {
		value = ToFloatValue(left.toFloat() + right.toFloat())
	} else {
		value = ToIntValue(left.toInt() + right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Subtract struct{}

func (self _Subtract) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isFloat() || right.isFloat() {
		value = ToFloatValue(left.toFloat() - right.toFloat())
	} else {
		value = ToIntValue(left.toInt() - right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Multiply struct{}

func (self _Multiply) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isFloat() || right.isFloat() {
		value = ToFloatValue(left.toFloat() * right.toFloat())
	} else {
		value = ToIntValue(left.toInt() * right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Divide struct{}

func (self _Divide) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isFloat() || right.isFloat() {
		value = ToFloatValue(left.toFloat() / right.toFloat())
	} else {
		value = ToIntValue(left.toInt() / right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}
