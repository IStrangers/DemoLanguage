package interpreter

import (
	"fmt"
	"github.com/istrangers/demolanguage/ast"
	"github.com/istrangers/demolanguage/token"
)

func (self *Interpreter) evaluateExpression(expression ast.Expression) Value {
	switch expr := expression.(type) {
	case *ast.NullLiteral:
		return self.evaluateNull()
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
	case *ast.ObjectLiteral:
		return self.evaluateObjectLiteral(expr)
	case *ast.ArrayLiteral:
		return self.evaluateArrayLiteral(expr)
	case *ast.FunLiteral:
		return self.evaluateFunLiteral(expr)
	case *ast.ArrowFunctionLiteral:
		return self.evaluateArrowFunctionLiteral(expr)
	case *ast.ThisExpression:
		return self.evaluateThisExpression(expr)
	case *ast.AssignExpression:
		return self.evaluateAssignExpression(expr)
	case *ast.BinaryExpression:
		return self.evaluateBinaryExpression(expr)
	case *ast.UnaryExpression:
		return self.evaluateUnaryExpression(expr)
	case *ast.CallExpression:
		return self.evaluateCallExpression(expr)
	case *ast.DotExpression:
		return self.evaluateDotExpression(expr)
	case *ast.BracketExpression:
		return self.evaluateBracketExpression(expr)
	}
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateSkip() Value {
	return Const_Skip_Value
}

func (self *Interpreter) evaluateBreak() Value {
	return Const_Break_Value
}

func (self *Interpreter) evaluateContinue() Value {
	return Const_Continue_Value
}

func (self *Interpreter) evaluateReturn(values []Value) Value {
	valueLength := len(values)
	if valueLength == 0 {
		return Const_Return_Value
	} else if valueLength == 1 {
		return Value{Return, values[0]}
	}
	return Value{Return, MultipleValueValue(values)}
}

func (self *Interpreter) evaluateNull() Value {
	return Const_Null_Value
}

func (self *Interpreter) evaluateBooleanLiteral(value bool) Value {
	if value {
		return Const_True_Value
	}
	return Const_False_Value
}

func (self *Interpreter) evaluateNumberLiteral(value any) Value {
	return NumberValue(value)
}

func (self *Interpreter) evaluateStringLiteral(value any) Value {
	return StringValue(value)
}

func (self *Interpreter) evaluateBinding(binding *ast.Binding) Value {
	targetValue := self.evaluateExpression(binding.Target)
	targetRef := targetValue.referenced()
	if self.runtime.getStash().contains(targetRef.getName()) {
		self.panic("already defined: "+targetRef.getName(), binding.StartIndex())
	}
	initValue := self.evaluateExpression(binding.Initializer)
	targetRef.setValue(initValue)
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateIdentifier(identifier *ast.Identifier) Value {
	name := identifier.Name
	return ReferenceValue(StashReferenced{
		name,
		self.runtime.getStash(),
	})
}

func (self *Interpreter) evaluateObjectLiteral(objectLiteral *ast.ObjectLiteral) Value {
	propertys := make(map[string]Value)
	for _, property := range objectLiteral.Properties {
		kv := property.(*ast.PropertyKeyValue)
		propertys[kv.Name.Name] = self.evaluateExpression(kv.Value)
	}
	return BuiltinObjectObject(propertys)
}

func (self *Interpreter) evaluateArrayLiteral(arrayLiteral *ast.ArrayLiteral) Value {
	var values []Value
	for _, value := range arrayLiteral.Values {
		values = append(values, self.evaluateExpression(value))
	}
	return BuiltinArrayObject(values)
}

func (self *Interpreter) evaluateFunction(identifier *ast.Identifier, parameterList *ast.ParameterList, body ast.Statement, funDefinition string) Value {
	globalFunction := Functiond{
		funDefinition: funDefinition,
		callee: func(arguments ...Value) Value {
			return self.evaluateCallFunction(parameterList, body, arguments...)
		},
	}
	functionValue := FunctionValue(globalFunction)

	if identifier == nil {
		return functionValue
	}
	globalFunction.name = identifier.Name
	identifierValue := self.evaluateExpression(identifier)
	identifierRef := identifierValue.referenced()
	identifierRef.setValue(functionValue)
	return self.evaluateSkip()
}

func (self *Interpreter) evaluateFunLiteral(funLiteral *ast.FunLiteral) Value {
	identifier := funLiteral.Name
	return self.evaluateFunction(identifier, funLiteral.ParameterList, funLiteral.Body, funLiteral.FunDefinition)
}

func (self *Interpreter) evaluateArrowFunctionLiteral(arrowFunctionLiteral *ast.ArrowFunctionLiteral) Value {
	return self.evaluateFunction(nil, arrowFunctionLiteral.ParameterList, arrowFunctionLiteral.Body, arrowFunctionLiteral.FunDefinition)
}

func (self *Interpreter) evaluateThisExpression(thisExpression *ast.ThisExpression) Value {
	return ObjectValue(self.runtime.scope.this)
}

func (self *Interpreter) evaluateAssignExpression(assignExpression *ast.AssignExpression) Value {
	leftValue := self.evaluateExpression(assignExpression.Left)
	rightValue := self.evaluateExpression(assignExpression.Right)
	if assignExpression.Operator != token.ASSIGN {
		rightValue = self.evaluateBinary(leftValue, assignExpression.Operator, rightValue)
	}
	leftRef := leftValue.referenced()
	leftRef.setValue(rightValue)
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
		if leftValue.isNumber() && rightValue.isNumber() {
			return self.evaluateBooleanLiteral(leftValue.float64() == rightValue.float64())
		}
		return self.evaluateBooleanLiteral(leftValue.getVal() == rightValue.getVal())
	case token.NOT_EQUAL:
		if leftValue.isNumber() && rightValue.isNumber() {
			return self.evaluateBooleanLiteral(leftValue.float64() != rightValue.float64())
		}
		return self.evaluateBooleanLiteral(leftValue.getVal() != rightValue.getVal())
	case token.LESS:
		return self.evaluateBooleanLiteral(leftValue.float64() < rightValue.float64())
	case token.LESS_OR_EQUAL:
		return self.evaluateBooleanLiteral(leftValue.float64() < rightValue.float64() || leftValue.float64() == rightValue.float64())
	case token.GREATER:
		return self.evaluateBooleanLiteral(leftValue.float64() > rightValue.float64())
	case token.GREATER_OR_EQUAL:
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
	return self.panic("Unsupported operator: "+operator.String(), -1)
}

func (self *Interpreter) evaluateUnaryExpression(unaryExpression *ast.UnaryExpression) Value {
	operandValue := self.evaluateExpression(unaryExpression.Operand)

	if operandValue.isReferenced() {
		operandRef := operandValue.referenced()
		switch unaryExpression.Operator {
		case token.NOT:
			operandRef.setValue(self.evaluateBooleanLiteral(!operandValue.bool()))
			return operandValue
		case token.ADDITION:
			val := operandValue.float64()
			if val < 0 {
				val = -val
			}
			operandRef.setValue(self.evaluateNumberLiteral(val))
			return operandValue
		case token.SUBTRACT:
			operandRef.setValue(self.evaluateNumberLiteral(-operandValue.float64()))
			return operandValue
		case token.INCREMENT:
			operandRef.setValue(self.evaluateNumberLiteral(operandValue.float64() + 1))
			return operandValue
		case token.DECREMENT:
			operandRef.setValue(self.evaluateNumberLiteral(operandValue.float64() - 1))
			return operandValue
		}
	} else {
		switch unaryExpression.Operator {
		case token.ADDITION:
			val := operandValue.float64()
			if val < 0 {
				val = -val
			}
			return self.evaluateNumberLiteral(val)
		case token.SUBTRACT:
			return self.evaluateNumberLiteral(-operandValue.float64())
		}
	}

	return self.evaluateSkip()
}

func (self *Interpreter) evaluateCallExpression(callExpression *ast.CallExpression) Value {
	calleeValue := self.evaluateExpression(callExpression.Callee)
	var calleeRef Referenced
	if calleeValue.getValueType() == Function {
		calleeRef = AnonymousReferenced{
			calleeValue,
		}
	} else {
		calleeRef = calleeValue.referenced()
	}
	calleeValue = calleeValue.flatResolve()
	if !calleeValue.isFunction() {
		self.panic(fmt.Sprintf("%s is not a function", calleeRef.getName()), callExpression.StartIndex())
	}
	function := calleeValue.functiond()
	var arguments []Value
	for _, argument := range callExpression.Arguments {
		arguments = append(arguments, self.evaluateExpression(argument))
	}
	var this Objectd
	var callee string
	if calleeRef.getType() == PropertyReferencedType {
		this = calleeRef.(PropertyReferenced).object
		callee = calleeRef.getName()
	} else {
		this = self.runtime.global
		callee = function.getName()
	}
	self.runtime.openScope(this, callee)
	defer self.runtime.closeScope()
	resultValue := function.call(arguments...)
	if resultValue.isReturn() {
		return resultValue.ofValue()
	}
	return resultValue
}

func (self *Interpreter) evaluateCallFunction(parameterList *ast.ParameterList, body ast.Statement, arguments ...Value) Value {
	argsLength := len(arguments)
	for index, binding := range parameterList.List {
		targetValue := self.evaluateExpression(binding.Target)
		targetRef := targetValue.referenced()
		if argsLength > index {
			targetRef.setValue(arguments[index])
		} else if binding.Initializer != nil {
			targetRef.setValue(self.evaluateExpression(binding.Initializer))
		}
	}
	return self.evaluateStatement(body)
}

func (self *Interpreter) evaluateDotExpression(dotExpression *ast.DotExpression) Value {
	leftValue := self.evaluateExpression(dotExpression.Left)
	leftValue = leftValue.flatResolve()
	leftObject := leftValue.objectd()
	identifier := dotExpression.Identifier
	return ReferenceValue(PropertyReferenced{
		identifier.Name,
		self.evaluateStringLiteral(identifier.Name),
		leftObject,
	})
}

func (self *Interpreter) evaluateBracketExpression(bracketExpression *ast.BracketExpression) Value {
	leftValue := self.evaluateExpression(bracketExpression.Left)
	leftValue = leftValue.flatResolve()
	leftObject := leftValue.objectd()
	exprValue := self.evaluateExpression(bracketExpression.Expression)
	return ReferenceValue(PropertyReferenced{
		exprValue.string(),
		exprValue,
		leftObject,
	})
}
