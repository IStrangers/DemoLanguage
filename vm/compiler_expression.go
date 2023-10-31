package vm

import "DemoLanguage/ast"

type CompiledExpression interface {
	isConstExpr() bool
}

type CompiledBaseExpression struct {
	compile *Compiler
	offset  int
}

type CompiledLiteralExpression struct {
	base  CompiledBaseExpression
	value Value
}

func (self CompiledLiteralExpression) isConstExpr() bool {
	return true
}

func (self *Compiler) compileExpression(expression ast.Expression) CompiledExpression {
	switch expr := expression.(type) {
	case *ast.NullLiteral:
		return self.compileNullLiteral(expr)
	case *ast.NumberLiteral:
		return self.compileNumberLiteral(expr)
	case *ast.StringLiteral:
		return self.compileStringLiteral(expr)
	default:
		return self.errorAssert(false, int(expression.StartIndex())-1, "Unknown expression type: %T", expression)
	}
}

func (self *Compiler) compileNullLiteral(expr *ast.NullLiteral) CompiledExpression {
	return &CompiledLiteralExpression{
		CompiledBaseExpression{self, int(expr.StartIndex()) - 1},
		Const_Null_Value,
	}
}

func (self *Compiler) compileNumberLiteral(expr *ast.NumberLiteral) CompiledExpression {
	var value Value
	switch val := expr.Value.(type) {
	case int64:
		value = ToIntValue(val)
	case float64:
		value = ToFloatValue(val)
	default:
		return self.errorAssert(false, int(expr.StartIndex())-1, "Unsupported number literal type: %T", expr.Value)
	}
	return &CompiledLiteralExpression{
		CompiledBaseExpression{self, int(expr.StartIndex()) - 1},
		value,
	}
}

func (self *Compiler) compileStringLiteral(expr *ast.StringLiteral) CompiledExpression {
	return &CompiledLiteralExpression{
		CompiledBaseExpression{self, int(expr.StartIndex()) - 1},
		ToStringValue(expr.Value),
	}
}
