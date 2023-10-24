package interpreter

import (
	"DemoLanguage/ast"
	"DemoLanguage/parser"
)

type Interpreter struct {
	runtime *Runtime
}

func CreateInterpreter() *Interpreter {
	return &Interpreter{
		runtime: createRunTime(),
	}
}

func (self *Interpreter) run(fileName string, content string) Value {
	parser := parser.CreateParser(1, fileName, content)
	program, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}
	return self.runProgram(program)
}

func (self *Interpreter) runProgram(program *ast.Program) Value {
	self.runtime.openScope()
	defer self.runtime.closeScope()
	value := self.evaluateProgramBody(program.Body)
	return value
}
