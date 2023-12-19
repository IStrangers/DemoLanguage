package interpreter

import (
	"fmt"
	"github.com/istrangers/demolanguage/ast"
	"github.com/istrangers/demolanguage/file"
	"github.com/istrangers/demolanguage/parser"
)

type Interpreter struct {
	runtime *Runtime
	file    *file.File
}

func CreateInterpreter() *Interpreter {
	return &Interpreter{
		runtime: createRunTime(),
	}
}

func (self *Interpreter) run(fileName string, content string) Value {
	parser := parser.CreateParser(1, fileName, content, true, true)
	program, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}
	self.file = program.File
	return self.runProgram(program)
}

func (self *Interpreter) runProgram(program *ast.Program) Value {
	self.runtime.openScope(self.runtime.global, "")
	defer self.runtime.closeScope()
	value := self.evaluateProgramBody(program.Body)
	return value
}

func (self *Interpreter) panic(msg string, index file.Index) Value {
	trace := ""
	if index != -1 {
		position := self.file.PositionByIndex(index)
		trace = fmt.Sprintf(" at %s (%s %d:%d)", self.runtime.scope.callee, position.FileName, position.Line, position.Column)
	}
	panic(fmt.Sprintf("%s\n\t%s", msg, trace))
	return Const_Skip_Value
}
