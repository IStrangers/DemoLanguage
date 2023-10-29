package vm

import (
	"testing"
)

func TestVM(t *testing.T) {
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
