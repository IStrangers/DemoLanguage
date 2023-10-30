package vm

import "DemoLanguage/file"

type Program struct {
	values       ValueArray
	instructions InstructionArray
	file         *file.File
}

type ValueStack ValueArray

func (self *ValueStack) expand(index int) {
	if index < len(*self) {
		return
	}
	index++
	var newCap int
	if index < 1024 {
		newCap = index * 2
	} else {
		newCap = (index + 1025) &^ 1023
	}
	newValueStack := make(ValueStack, index, newCap)
	copy(newValueStack, *self)
	*self = newValueStack
}

type VM struct {
	program *Program
	stack   ValueStack
	pc      int
	sp      int
}

func (self *VM) run() {
	for {
		if self.pc >= self.getInstructionSize() {
			break
		}
		self.execInstruction(self.pc)
	}
}

func (self *VM) getInstructionSize() int {
	return len(self.program.instructions)
}

func (self *VM) getInstruction(pc int) Instruction {
	return self.program.instructions[pc]
}

func (self *VM) execInstruction(pc int) {
	instruction := self.getInstruction(pc)
	instruction.exec(self)
}

func (self *VM) getValue(index uint32) Value {
	return self.program.values[index]
}

func (self *VM) push(value Value) {
	self.stack.expand(self.sp)
	self.stack[self.sp] = value
	self.sp++
}

func (self *VM) pop() Value {
	self.sp--
	return self.stack[self.sp]
}

func (self *VM) peek() Value {
	return self.stack[self.sp-1]
}

func (self *VM) clearStack() {
	sp := self.sp
	stackTail := self.stack[sp:]
	for i := range stackTail {
		stackTail[i] = nil
	}
	self.stack = self.stack[:sp]
}
