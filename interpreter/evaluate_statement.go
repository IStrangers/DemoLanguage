package interpreter

import (
	"DemoLanguage/ast"
	"DemoLanguage/token"
)

func (self *Interpreter) evaluateProgramBody(listStatement []ast.Statement) Value {
	var value Value
	var isResult bool
	for _, st := range listStatement {
		value = self.evaluateStatement(st)
		isResult = value.isResult()
		if isResult && value.isReturn() {
			return value.ofValue()
		}
	}
	if isResult {
		return value
	}
	return self.evaluateNull()
}

func (self *Interpreter) evaluateStatement(statement ast.Statement) Value {
	switch st := statement.(type) {
	case *ast.BlockStatement:
		return self.evaluateBlockStatement(st)
	case *ast.VarStatement:
		return self.evaluateVarStatement(st)
	case *ast.BreakStatement:
		return self.evaluateBreakStatement(st)
	case *ast.ContinueStatement:
		return self.evaluateContinueStatement(st)
	case *ast.ReturnStatement:
		return self.evaluateReturnStatement(st)
	case *ast.IfStatement:
		return self.evaluateIfStatement(st)
	case *ast.SwitchStatement:
		return self.evaluateSwitchStatement(st)
	case *ast.ForStatement:
		return self.evaluateForStatement(st)
	case *ast.FunStatement:
		return self.evaluateFunStatement(st)
	case *ast.ExpressionStatement:
		return self.evaluateExpressionStatement(st)
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateBlockStatement(blockStatement *ast.BlockStatement) Value {
	for _, st := range blockStatement.Body {
		value := self.evaluateStatement(st)
		if value.isBreak() || value.isContinue() || value.isReturn() {
			return value
		}
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateVarStatement(varStatement *ast.VarStatement) Value {
	for _, binding := range varStatement.List {
		self.evaluateBinding(binding)
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateBreakStatement(breakStatement *ast.BreakStatement) Value {
	return self.evaluateBreak()
}

func (self *Interpreter) evaluateContinueStatement(continueStatement *ast.ContinueStatement) Value {
	return self.evaluateContinue()
}

func (self *Interpreter) evaluateReturnStatement(returnStatement *ast.ReturnStatement) Value {
	var values []Value
	for _, argument := range returnStatement.Arguments {
		values = append(values, self.evaluateExpression(argument))
	}
	return self.evaluateReturn(values)
}

func (self *Interpreter) evaluateIfStatement(ifStatement *ast.IfStatement) Value {
	conditionValue := self.evaluateExpression(ifStatement.Condition)
	if conditionValue.bool() {
		return self.evaluateStatement(ifStatement.Consequent)
	} else if ifStatement.Alternate != nil {
		return self.evaluateStatement(ifStatement.Alternate)
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateSwitchStatement(switchStatement *ast.SwitchStatement) Value {
	discriminantValue := self.evaluateExpression(switchStatement.Discriminant)
	var consequent ast.Statement
	for _, caseStatement := range switchStatement.Body {
		if caseStatement.Condition == nil {
			consequent = caseStatement.Consequent
			continue
		}
		comparisonValue := self.evaluateComparison(discriminantValue, token.EQUAL, self.evaluateExpression(caseStatement.Condition))
		if comparisonValue.bool() {
			consequent = caseStatement.Consequent
			break
		}
	}
	if consequent == nil {
		return self.evaluateSkip()
	}
	return self.evaluateStatement(consequent)
}

func (self *Interpreter) evaluateForStatement(forStatement *ast.ForStatement) Value {
	if forStatement.Initializer != nil {
		self.evaluateStatement(forStatement.Initializer)
	}
	for {
		if forStatement.Condition != nil {
			conditionValue := self.evaluateExpression(forStatement.Condition)
			if !conditionValue.bool() {
				break
			}
		}
		value := self.evaluateStatement(forStatement.Body)
		if value.isBreak() {
			break
		} else if value.isContinue() {
			// none
		} else if !value.isSkip() {
			return value
		}
		if forStatement.Update != nil {
			self.evaluateExpression(forStatement.Update)
		}
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateFunStatement(funStatement *ast.FunStatement) Value {
	return self.evaluateExpression(funStatement.FunLiteral)
}

func (self *Interpreter) evaluateExpressionStatement(expressionStatement *ast.ExpressionStatement) Value {
	return self.evaluateExpression(expressionStatement.Expression)
}
