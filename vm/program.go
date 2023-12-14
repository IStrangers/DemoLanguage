package vm

import (
	"DemoLanguage/file"
	"fmt"
)

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

func (self *Program) dumpInstructions(logger func(format string, args ...interface{})) {
	self.dumpInstructionsByIndent("", logger)
}

const (
	colorGreen = "\033[32m"
)

func colorize(color, message string) string {
	return fmt.Sprintf("%s%s", color, message)
}

func (self *Program) dumpInstructionsByIndent(indent string, logger func(format string, args ...interface{})) {
	logger(colorize(colorGreen, "values: %+v"), self.values)
	//dumpInitFields := func(initFields *Program) {
	//	i := indent + ">"
	//	logger("%s ---- init_fields:", i)
	//	initFields.dumpInstructionsByIndent(i, logger)
	//	logger("%s ----", i)
	//}
	for pc, ins := range self.instructions {
		logger(colorize(colorGreen, "%s %d: %T(%v)"), indent, pc, ins, ins)
		var prg *Program
		switch f := ins.(type) {
		case *NewFun:
			prg = f.program
		}
		if prg != nil {
			prg.dumpInstructionsByIndent(indent+">", logger)
		}
	}
}
