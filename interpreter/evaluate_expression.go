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
	return Value{Skip, nil}
}

func (self *Interpreter) evaluateNullLiteral() Value {
	return Value{Null, nil}
}

func (self *Interpreter) evaluateBooleanLiteral(value any) Value {
	return Value{Boolean, value}
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

func (self *Interpreter) evaluateReference(value any) Value {
	return Value{Reference, value}
}

func (self *Interpreter) evaluateBinding(binding *ast.Binding) Value {
	targetValue := self.evaluateExpression(binding.Target)
	targetRef := targetValue.reference()
	value := targetRef.getValue()
	if value.getValue() != nil {
		panic("already defined: " + targetRef.getName())
	}
	initValue := self.evaluateExpression(binding.Initializer)
	targetRef.setValue(initValue)
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateIdentifier(identifier *ast.Identifier) Value {
	name := identifier.Name
	return self.evaluateReference(&StashReference{
		name,
		self.runtime.getStash(),
	})
}

func (self *Interpreter) evaluateAssignExpression(assignExpression *ast.AssignExpression) Value {
	leftValue := self.evaluateExpression(assignExpression.Left)
	leftRef := leftValue.reference()
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
		return self.evaluateBooleanLiteral(leftValue.getValue() == rightValue.getValue())
	case token.NOT_EQUAL:
		return self.evaluateBooleanLiteral(leftValue.getValue() != rightValue.getValue())
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
	operandRef := operandValue.reference()
	switch unaryExpression.Operator {
	case token.NOT:
		operandValue.value = !operandValue.bool()
		operandRef.setValue(operandValue)
		return operandValue
	case token.INCREMENT:
		operandValue.value = operandValue.float64() + 1
		operandRef.setValue(operandValue)
		return operandValue
	case token.DECREMENT:
		operandValue.value = operandValue.float64() - 1
		operandRef.setValue(operandValue)
		return operandValue
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateCallExpression(callExpression *ast.CallExpression) Value {
	return self.evaluateSkip()
}
