package vm

import "strings"

type Instruction interface {
	exec(*VM)
}

type InstructionArray []Instruction

var (
	Addition _Addition
	Subtract _Subtract
	Multiply _Multiply
	Divide   _Divide

	Less           _Less
	LessOrEqual    _LessOrEqual
	Greater        _Greater
	GreaterOrEqual _GreaterOrEqual
)

type LoadValue uint32

func (self LoadValue) exec(vm *VM) {
	vm.push(vm.getValue(uint32(self)))
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

type _Less struct{}

func (self _Less) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := lessComp(left, right)

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _LessOrEqual struct{}

func (self _LessOrEqual) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_True_Value
	if lessComp(right, left) == Const_Bool_True_Value {
		value = Const_Bool_False_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Greater struct{}

func (self _Greater) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := lessComp(right, left)

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _GreaterOrEqual struct{}

func (self _GreaterOrEqual) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_True_Value
	if lessComp(left, right) == Const_Bool_True_Value {
		value = Const_Bool_False_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

func lessComp(left Value, right Value) Value {
	var less bool
	if left.isString() && right.isString() {
		less = strings.Compare(left.toString(), right.toString()) < 0
	} else if left.isInt() || left.isFloat() && right.isInt() || right.isFloat() {
		if left.isFloat() || right.isFloat() {
			less = left.toFloat() < right.toFloat()
		} else {
			less = left.toInt() < right.toInt()
		}
	}
	if less {
		return Const_Bool_True_Value
	}
	return Const_Bool_False_Value
}

type Jump int

func (self Jump) exec(vm *VM) {
	vm.pc += int(self)
}
