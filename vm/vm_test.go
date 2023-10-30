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
				LoadValue(0),
				LoadValue(1),
				Addition,
				LoadValue(1),
				Subtract,
				LoadValue(1),
				Multiply,
				LoadValue(1),
				Divide,
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
		LoadValue(0),
		LoadValue(1),
		Less,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
	vm.program.instructions = InstructionArray{
		LoadValue(0),
		LoadValue(1),
		LessOrEqual,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
	vm.program.instructions = InstructionArray{
		LoadValue(0),
		LoadValue(1),
		Greater,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
	vm.program.instructions = InstructionArray{
		LoadValue(0),
		LoadValue(1),
		GreaterOrEqual,
	}
	vm.run()
	println(vm.pop().toString())
	vm.pc = 0
	vm.clearStack()
}
