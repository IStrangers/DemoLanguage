package vm

import (
	"math"
	"strings"
)

type Instruction interface {
	exec(*VM)
}

type InstructionArray []Instruction

func (self *InstructionArray) size() int {
	return len(*self)
}

func (self *InstructionArray) add(instructions ...Instruction) {
	*self = append(*self, instructions...)
}

var (
	Add _Add
	Sub _Sub
	Mul _Mul
	Div _Div
	Mod _Mod

	AND _AND
	OR  _OR
	Not _Not
	Inc _Inc
	Dec _Dec
	Neg _Neg

	EQ _EQ
	NE _NE
	LT _LT
	LE _LE
	GT _GT
	GE _GE

	Pop     _Pop
	InitVar _InitVar
)

type LoadVal int

func (self LoadVal) exec(vm *VM) {
	vm.push(vm.getValue(int(self)))
	vm.pc++
}

type _Add struct{}

func (self _Add) exec(vm *VM) {
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

type _Sub struct{}

func (self _Sub) exec(vm *VM) {
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

type _Mul struct{}

func (self _Mul) exec(vm *VM) {
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

type _Div struct{}

func (self _Div) exec(vm *VM) {
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

type _Mod struct{}

func (self _Mod) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isFloat() || right.isFloat() {
		value = ToFloatValue(math.Mod(left.toFloat(), right.toFloat()))
	} else {
		value = ToIntValue(left.toInt() % right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _AND struct{}

func (self _AND) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := ToIntValue(left.toInt() & right.toInt())

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _OR struct{}

func (self _OR) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := ToIntValue(left.toInt() | right.toInt())

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Not struct{}

func (self _Not) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.toBool() {
		value = Const_Bool_False_Value
	} else {
		value = Const_Bool_True_Value
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _Inc struct{}

func (self _Inc) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.isFloat() {
		value = ToFloatValue(value.toFloat() + 1.0)
	} else {
		value = ToIntValue(value.toInt() + 1)
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _Dec struct{}

func (self _Dec) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.isFloat() {
		value = ToFloatValue(value.toFloat() - 1.0)
	} else {
		value = ToIntValue(value.toInt() - 1)
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _Neg struct{}

func (self _Neg) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.isFloat() {
		value = ToFloatValue(-value.toFloat())
	} else {
		value = ToIntValue(-value.toInt())
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _EQ struct{}

func (self _EQ) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_False_Value
	if left.equals(right) {
		value = Const_Bool_True_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _NE struct{}

func (self _NE) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_True_Value
	if left.equals(right) {
		value = Const_Bool_False_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _LT struct{}

func (self _LT) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := lessComp(left, right)

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _LE struct{}

func (self _LE) exec(vm *VM) {
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

type _GT struct{}

func (self _GT) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := lessComp(right, left)

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _GE struct{}

func (self _GE) exec(vm *VM) {
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

type Jeq int

func (self Jeq) exec(vm *VM) {
	vm.sp--
	value := vm.stack[vm.sp]
	if value.toBool() {
		vm.pc += int(self)
	} else {
		vm.pc++
	}
}

type Jne int

func (self Jne) exec(vm *VM) {
	vm.sp--
	value := vm.stack[vm.sp]
	if !value.toBool() {
		vm.pc += int(self)
	} else {
		vm.pc++
	}
}

type Jump int

func (self Jump) exec(vm *VM) {
	vm.pc += int(self)
}

type _Pop struct{}

func (self _Pop) exec(vm *VM) {
	vm.sp--
	vm.pc++
}

type ResolveVar string

func (self ResolveVar) exec(vm *VM) {
	vm.refStack.add(&ObjectRef{
		refObject: vm.runtime.globalObject,
		refName:   string(self),
	})
	vm.pc++
}

type _InitVar struct{}

func (self _InitVar) exec(vm *VM) {
	ref := vm.refStack.pop()
	vm.sp--
	ref.set(vm.stack[vm.sp])
	vm.pc++
}

type LoadVar string

func (self LoadVar) exec(vm *VM) {
	globalObject := vm.runtime.globalObject
	name := string(self)
	value := globalObject.self.getProperty(name)
	if value == nil {

	}
	vm.push(value)
	vm.pc++
}
