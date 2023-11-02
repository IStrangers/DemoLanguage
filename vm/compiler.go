package vm

import (
	"DemoLanguage/ast"
	"DemoLanguage/file"
	"fmt"
)

type CompilerError struct {
	Message string
	File    *file.File
	Offset  int
}

type CompilerSyntaxError struct {
	CompilerError
}

type CompilerReferenceError struct {
	CompilerError
}

type Compiler struct {
	program *Program
	scope   *Scope
	block   *Block
	evalVM  *VM
}

func (self *Compiler) compile(in *ast.Program) {
	body := in.Body
	self.compileStatements(body, true)
}

func (self *Compiler) addProgramValue(value Value) int {
	return self.program.addValue(value)
}

func (self *Compiler) addProgramInstructions(instructions ...Instruction) {
	self.program.addInstructions(instructions...)
}

func (self *Compiler) setProgramInstruction(index int, instruction Instruction) {
	self.program.setProgramInstruction(index, instruction)
}

func (self *Compiler) openScope() {
	self.scope = &Scope{
		outer:   self.scope,
		program: self.program,
	}
}

func (self *Compiler) closeScope() {
	self.scope = self.scope.outer
}

func (self *Compiler) enterVirtualMode() func() {
	originProgram, originBlock := self.program, self.block
	if originBlock != nil {
		self.block = &Block{
			originBlock.outer,
			originBlock.blockType,
			originBlock.label,
		}
	}
	self.program = &Program{
		source: self.program.source,
	}
	self.openScope()
	return func() {
		self.program, self.block = originProgram, originBlock
		self.closeScope()
	}
}

func (self *Compiler) throwSyntaxError(offset int, format string, args ...any) CompiledExpression {
	panic(&CompilerSyntaxError{
		CompilerError{
			File:    self.program.source,
			Offset:  offset,
			Message: fmt.Sprintf(format, args...),
		},
	})
	return nil
}

func (self *Compiler) errorAssert(cond bool, offset int, message string, args ...any) CompiledExpression {
	if cond {
		return nil
	}
	return self.throwSyntaxError(offset, "Compiler error: "+message, args...)
}
