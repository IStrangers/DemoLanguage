package vm

import "DemoLanguage/ast"

type Compiler struct {
	program *Program
}

func (self *Compiler) compile(in *ast.Program) {
	body := in.Body
	self.compileStatements(body, true)
}
