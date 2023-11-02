package vm

import "DemoLanguage/ast"

func (self *Compiler) compileStatements(statements []ast.Statement, needResult bool) {
	for _, statement := range statements {
		self.compileStatement(statement, needResult)
	}
}

func (self *Compiler) compileStatement(statement ast.Statement, needResult bool) {
	switch st := statement.(type) {
	case *ast.BlockStatement:
		self.compileBlockStatement(st, needResult)
	case *ast.IfStatement:
		self.compileIfStatement(st, needResult)
	case *ast.ExpressionStatement:
		self.compileExpressionStatement(st, needResult)
	}
}

func (self *Compiler) compileBlockStatement(st *ast.BlockStatement, needResult bool) {
	self.compileStatements(st.Body, needResult)
}

func (self *Compiler) compileIfStatement(st *ast.IfStatement, needResult bool) {
	conditionExpr := self.compileExpression(st.Condition)
	if conditionExpr.isConstExpression() {
		res, ex := self.evalConstExpr(conditionExpr)
		if ex != nil {
			conditionExpr.addSourceMap()
			self.emitThrow(ex.value)
			return
		}
		if res.toBool() {
			self.compileStatement(st.Consequent, needResult)
		} else if st.Alternate != nil {
			self.compileStatement(st.Alternate, needResult)
		}
	} else {

	}
}

func (self *Compiler) compileExpressionStatement(st *ast.ExpressionStatement, needResult bool) {
	self.compileExpression(st.Expression)
}
