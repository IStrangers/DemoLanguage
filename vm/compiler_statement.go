package vm

import "DemoLanguage/ast"

func (self *Compiler) compileStatements(statements []ast.Statement, needResult bool) {
	for _, statement := range statements {
		self.compileStatement(statement, needResult)
	}
}

func (self *Compiler) compileStatement(statement ast.Statement, needResult bool) {
	switch st := statement.(type) {
	case *ast.IfStatement:
		self.compileIfStatement(st, needResult)
	case *ast.ExpressionStatement:
		self.compileExpressionStatement(st, needResult)
	}
}

func (self *Compiler) compileIfStatement(st *ast.IfStatement, needResult bool) {
	conditionExpr := self.compileExpression(st.Condition)
	if conditionExpr.isConstExpr() {

	} else {

	}
}

func (self *Compiler) compileExpressionStatement(st *ast.ExpressionStatement, needResult bool) {
	self.compileExpression(st.Expression)
}
