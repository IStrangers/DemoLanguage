package vm

import "DemoLanguage/ast"

type CompiledExpression interface {
}

func (self *Compiler) compileExpression(expression ast.Expression) CompiledExpression {
	switch expr := expression.(type) {
	case *ast.NullLiteral:
		return self.compileNullLiteral(expr)
	case *ast.NumberLiteral:
		return self.compileNumberLiteral(expr)
	case *ast.StringLiteral:
		return self.compileStringLiteral(expr)
	}
}

func (self *Compiler) compileNullLiteral(expr *ast.NullLiteral) CompiledExpression {

}

func (self *Compiler) compileNumberLiteral(expr *ast.NumberLiteral) CompiledExpression {

}

func (self *Compiler) compileStringLiteral(expr *ast.StringLiteral) CompiledExpression {

}
