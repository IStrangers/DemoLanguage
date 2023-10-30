package vm

import (
	"testing"
)

func TestOperator(t *testing.T) {
	vm := &VM{
		pc: 0,
		program: &Program{
			values: ValueArray{ToFloatValue(5.99), ToIntValue(10)},
			instructions: InstructionArray{
				LoadVal(0),
				LoadVal(1),
				Add,
				LoadVal(1),
				Sub,
				LoadVal(1),
				Mul,
				LoadVal(1),
				Div,
				LoadVal(1),
				Mod,
			},
		},
	}
	vm.run()
	value := vm.pop()
	println(value.toString())
}

func TestLogicalOperator(t *testing.T) {
	vm := &VM{
		pc: 0,
		program: &Program{
			values: ValueArray{ToFloatValue(5.99), ToIntValue(10)},
		},
	}
	vm.program.instructions = InstructionArray{
		LoadVal(0),
		LoadVal(1),
		LT,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
	vm.program.instructions = InstructionArray{
		LoadVal(0),
		LoadVal(1),
		LE,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
	vm.program.instructions = InstructionArray{
		LoadVal(0),
		LoadVal(1),
		GT,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
	vm.program.instructions = InstructionArray{
		LoadVal(0),
		LoadVal(1),
		GE,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
}
