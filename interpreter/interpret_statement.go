package interpreter

import "DemoLanguage/ast"

func (self *Interpreter) evaluateListStatement(listStatement []ast.Statement) Value {
	for _, st := range listStatement {
		value := self.evaluateStatement(st)
		if !value.isNIL() {
			return value
		}
	}
	return self.evaluateNIL()
}

func (self *Interpreter) evaluateStatement(statement ast.Statement) Value {
	switch st := statement.(type) {
	case *ast.BlockStatement:
		return self.evaluateBlockStatement(st)
	case *ast.ReturnStatement:
		return self.evaluateReturnStatement(st)
	case *ast.IfStatement:
		return self.evaluateIfStatement(st)
	case *ast.ExpressionStatement:
		return self.evaluateExpressionStatement(st)
	}
	return self.evaluateNIL()
}

func (self *Interpreter) evaluateBlockStatement(blockStatement *ast.BlockStatement) Value {
	return self.evaluateListStatement(blockStatement.Body)
}

func (self *Interpreter) evaluateReturnStatement(returnStatement *ast.ReturnStatement) Value {
	var values []Value
	for _, argument := range returnStatement.Arguments {
		values = append(values, self.evaluateExpression(argument))
	}
	return self.evaluateObject(values)
}

func (self *Interpreter) evaluateIfStatement(ifStatement *ast.IfStatement) Value {
	conditionValue := self.evaluateExpression(ifStatement.Condition)
	if conditionValue.bool() {
		return self.evaluateStatement(ifStatement.Consequent)
	} else if ifStatement.Alternate != nil {
		return self.evaluateStatement(ifStatement.Alternate)
	}
	return self.evaluateNIL()
}

func (self *Interpreter) evaluateExpressionStatement(expressionStatement *ast.ExpressionStatement) Value {
	return self.evaluateExpression(expressionStatement.Expression)
}
