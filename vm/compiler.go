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

func (self CompilerSyntaxError) Error() string {
	if self.File != nil {
		return fmt.Sprintf("SyntaxError: %s at %s", self.Message, self.File.Position(self.Offset))
	}
	return fmt.Sprintf("SyntaxError: %s", self.Message)
}

type CompilerReferenceError struct {
	CompilerError
}

type Compiler struct {
	program *Program
	scope   *Scope
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
	originProgram := self.program
	self.program = &Program{
		source: self.program.source,
	}
	self.openScope()
	return func() {
		self.program = originProgram
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
