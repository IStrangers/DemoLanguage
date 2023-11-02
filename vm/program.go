package vm

import "DemoLanguage/file"

type SourceMapItem struct {
	pc  int
	pos int
}

type SourceMapItemArray []SourceMapItem

func (self *SourceMapItemArray) add(item SourceMapItem) {
	*self = append(*self, item)
}

func (self *SourceMapItemArray) size() int {
	return len(*self)
}

type Program struct {
	values       ValueArray
	instructions InstructionArray
	functionName string
	source       *file.File
	sourceMaps   SourceMapItemArray
}

func (self *Program) addValue(value Value) int {
	index := self.values.findIndex(value)
	if index != -1 {
		return index
	}
	self.values.add(value)
	return self.values.size() - 1
}

func (self *Program) getInstructionSize() int {
	return self.instructions.size()
}

func (self *Program) addInstructions(instructions ...Instruction) {
	self.instructions.add(instructions...)
}

func (self *Program) setProgramInstruction(index int, instruction Instruction) {
	self.instructions[index] = instruction
}

func (self *Program) getInstruction(index int) Instruction {
	return self.instructions[index]
}

func (self *Program) addSourceMap(pos int) {
	if len(self.sourceMaps) > 0 && self.sourceMaps[self.sourceMaps.size()-1].pos == pos {
		return
	}
	self.sourceMaps.add(SourceMapItem{int(self.instructions.size()), pos})
}
