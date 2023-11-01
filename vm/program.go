package vm

import "DemoLanguage/file"

type Program struct {
	values       ValueArray
	instructions InstructionArray
	file         *file.File
}

func (self *Program) addValue(value Value) uint {
	index := self.values.findIndex(value)
	if index != -1 {
		return uint(index)
	}
	self.values.add(value)
	return uint(self.values.size() - 1)
}

func (self *Program) getInstructionSize() int {
	return self.instructions.size()
}

func (self *Program) addInstructions(instructions ...Instruction) {
	self.instructions.add(instructions...)
}

func (self *Program) getInstruction(index int) Instruction {
	return self.instructions[index]
}
