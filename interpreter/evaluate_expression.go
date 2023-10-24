package interpreter

import (
	"DemoLanguage/ast"
	"DemoLanguage/token"
)

func (self *Interpreter) evaluateExpression(expression ast.Expression) Value {
	switch expr := expression.(type) {
	case *ast.NullLiteral:
		return self.evaluateNullLiteral()
	case *ast.BooleanLiteral:
		return self.evaluateBooleanLiteral(expr.Value)
	case *ast.NumberLiteral:
		return self.evaluateNumberLiteral(expr.Value)
	case *ast.StringLiteral:
		return self.evaluateStringLiteral(expr.Value)
	case *ast.Binding:
		return self.evaluateBinding(expr)
	case *ast.Identifier:
		return self.evaluateIdentifier(expr)
	case *ast.FunLiteral:
		return self.evaluateFunLiteral(expr)
	case *ast.AssignExpression:
		return self.evaluateAssignExpression(expr)
	case *ast.BinaryExpression:
		return self.evaluateBinaryExpression(expr)
	case *ast.UnaryExpression:
		return self.evaluateUnaryExpression(expr)
	case *ast.CallExpression:
		return self.evaluateCallExpression(expr)
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateSkip() Value {
	return self.evaluateSkipValue(nil)
}
func (self *Interpreter) evaluateSkipValue(value any) Value {
	return Value{Skip, value}
}

func (self *Interpreter) evaluateNullLiteral() Value {
	return Value{Null, nil}
}

func (self *Interpreter) evaluateBooleanLiteral(value bool) Value {
	if value {
		return Const_True_Value
	}
	return Const_False_Value
}

func (self *Interpreter) evaluateNumberLiteral(value any) Value {
	return Value{Number, value}
}

func (self *Interpreter) evaluateStringLiteral(value any) Value {
	return Value{String, value}
}

func (self *Interpreter) evaluateObject(value any) Value {
	return Value{Object, value}
}

func (self *Interpreter) evaluateFunction(value any) Value {
	return Value{Function, value}
}

func (self *Interpreter) evaluateReference(value any) Value {
	return Value{Reference, value}
}

func (self *Interpreter) evaluateBinding(binding *ast.Binding) Value {
	targetValue := self.evaluateExpression(binding.Target)
	targetRef := targetValue.referenced()
	if targetRef.getVal() != nil {
		panic("already defined: " + targetRef.getName())
	}
	initValue := self.evaluateExpression(binding.Initializer)
	targetRef.setValue(initValue)
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateIdentifier(identifier *ast.Identifier) Value {
	name := identifier.Name
	return self.evaluateReference(&StashReferenced{
		name,
		self.runtime.getStash(),
	})
}

func (self *Interpreter) evaluateFunLiteral(funLiteral *ast.FunLiteral) Value {
	identifier := funLiteral.Name
	identifierValue := self.evaluateExpression(identifier)
	identifierRef := identifierValue.referenced()
	globalFunctiond := &GlobalFunctiond{
		name: identifier.Name,
		callee: func(arguments ...Value) Value {
			return self.evaluateCallFunction(funLiteral, arguments...)
		},
	}
	identifierRef.setValue(self.evaluateFunction(globalFunctiond))
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateAssignExpression(assignExpression *ast.AssignExpression) Value {
	leftValue := self.evaluateExpression(assignExpression.Left)
	leftRef := leftValue.referenced()
	leftRef.setValue(self.evaluateExpression(assignExpression.Right))
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateBinaryExpression(binaryExpression *ast.BinaryExpression) Value {
	left, operator, right, comparison := binaryExpression.Left, binaryExpression.Operator, binaryExpression.Right, binaryExpression.Comparison

	leftValue := self.evaluateExpression(left)

	if operator == token.LOGICAL_OR {
		if leftValue.bool() {
			return leftValue
		}
		return self.evaluateExpression(right)
	}

	rightValue := self.evaluateExpression(right)

	if operator == token.LOGICAL_AND {
		if !leftValue.bool() {
			return leftValue
		}
		return rightValue
	}

	if comparison {
		return self.evaluateComparison(leftValue, operator, rightValue)
	} else {
		return self.evaluateBinary(leftValue, operator, rightValue)
	}
}

func (self *Interpreter) evaluateComparison(leftValue Value, operator token.Token, rightValue Value) Value {
	switch operator {
	case token.EQUAL:
		if leftValue.isReferenceNumber() && rightValue.isReferenceNumber() {
			return self.evaluateBooleanLiteral(leftValue.float64() == rightValue.float64())
		}
		return self.evaluateBooleanLiteral(leftValue.getVal() == rightValue.getVal())
	case token.NOT_EQUAL:
		if leftValue.isReferenceNumber() && rightValue.isReferenceNumber() {
			return self.evaluateBooleanLiteral(leftValue.float64() != rightValue.float64())
		}
		return self.evaluateBooleanLiteral(leftValue.getVal() != rightValue.getVal())
	case token.LESS:
		return self.evaluateBooleanLiteral(leftValue.float64() < rightValue.float64())
	case token.LESS_OR_EQUAL:
		return self.evaluateBooleanLiteral(leftValue.float64() < rightValue.float64() || leftValue.float64() == rightValue.float64())
	case token.GREATER:
		return self.evaluateBooleanLiteral(leftValue.float64() > rightValue.float64())
	case token.GREATER_OR_EQUEAL:
		return self.evaluateBooleanLiteral(leftValue.float64() > rightValue.float64() || leftValue.float64() == rightValue.float64())
	}
	return self.evaluateBooleanLiteral(false)
}

func (self *Interpreter) evaluateBinary(leftValue Value, operator token.Token, rightValue Value) Value {
	switch operator {
	case token.ADDITION:
		if leftValue.isString() || rightValue.isString() {
			return self.evaluateStringLiteral(leftValue.string() + rightValue.string())
		} else {
			return self.evaluateNumberLiteral(leftValue.float64() + rightValue.float64())
		}
	case token.SUBTRACT:
		return self.evaluateNumberLiteral(leftValue.float64() - rightValue.float64())
	case token.MULTIPLY:
		return self.evaluateNumberLiteral(leftValue.float64() * rightValue.float64())
	case token.DIVIDE:
		return self.evaluateNumberLiteral(leftValue.float64() / rightValue.float64())
	case token.REMAINDER:
		return self.evaluateNumberLiteral(leftValue.int64() % rightValue.int64())
	case token.AND_ARITHMETIC:
		return self.evaluateNumberLiteral(leftValue.int64() & rightValue.int64())
	case token.OR_ARITHMETIC:
		return self.evaluateNumberLiteral(leftValue.int64() | rightValue.int64())
	}
	panic("Unsupported operator: " + operator.String())
}

func (self *Interpreter) evaluateUnaryExpression(unaryExpression *ast.UnaryExpression) Value {
	operandValue := self.evaluateExpression(unaryExpression.Operand)
	operandRef := operandValue.referenced()
	switch unaryExpression.Operator {
	case token.NOT:
		operandRef.setValue(self.evaluateBooleanLiteral(!operandValue.bool()))
		return operandValue
	case token.INCREMENT:
		operandRef.setValue(self.evaluateNumberLiteral(operandValue.float64() + 1))
		return operandValue
	case token.DECREMENT:
		operandRef.setValue(self.evaluateNumberLiteral(operandValue.float64() - 1))
		return operandValue
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateCallExpression(callExpression *ast.CallExpression) Value {
	calleeValue := self.evaluateExpression(callExpression.Callee)
	calleeRef := calleeValue.referenced()
	functiondValue := calleeRef.getValue()
	functiond := functiondValue.functiond()
	var arguments []Value
	for _, argument := range callExpression.Arguments {
		arguments = append(arguments, self.evaluateExpression(argument))
	}
	return functiond.call(arguments...)
}

func (self *Interpreter) evaluateCallFunction(funLiteral *ast.FunLiteral, arguments ...Value) Value {
	self.runtime.openScope()
	defer self.runtime.closeScope()
	argsLength := len(arguments)
	for index, binding := range funLiteral.ParameterList.List {
		targetValue := self.evaluateExpression(binding.Target)
		targetRef := targetValue.referenced()
		if argsLength > index {
			targetRef.setValue(arguments[index])
		} else if binding.Initializer != nil {
			targetRef.setValue(self.evaluateExpression(binding.Initializer))
		}
	}
	return self.evaluateStatement(funLiteral.Body)
}
