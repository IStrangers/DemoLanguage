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
}

func (self *Compiler) compile(in *ast.Program) {
	body := in.Body
	self.compileStatements(body, true)
}

func (self *Compiler) throwSyntaxError(offset int, format string, args ...any) CompiledExpression {
	panic(&CompilerSyntaxError{
		CompilerError{
			File:    self.program.file,
			Offset:  offset,
			Message: fmt.Sprintf(format, args),
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
