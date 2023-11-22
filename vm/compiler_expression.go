package vm

import (
	"DemoLanguage/ast"
	"DemoLanguage/file"
	"DemoLanguage/token"
)

type CompiledExpression interface {
	isLiteralExpression() bool
	isConstExpression() bool
	addSourceMap()
}

type CompiledBaseExpression struct {
	compile *Compiler
	offset  int
}

func (self CompiledBaseExpression) isLiteralExpression() bool {
	return false
}

func (self CompiledBaseExpression) addSourceMap() {
	if self.offset < 0 {
		return
	}
	self.compile.program.addSourceMap(self.offset)
}

type CompiledLiteralExpression struct {
	CompiledBaseExpression
	value Value
}

func (self CompiledLiteralExpression) isLiteralExpression() bool {
	return true
}

func (self CompiledLiteralExpression) isConstExpression() bool {
	return true
}

type CompiledIdentifierExpression struct {
	CompiledBaseExpression
	name string
}

func (self CompiledIdentifierExpression) isConstExpression() bool {
	return false
}

type CompiledUnaryExpression struct {
	CompiledBaseExpression
	operator token.Token
	operand  CompiledExpression
	postfix  bool
}

func (self CompiledUnaryExpression) isConstExpression() bool {
	return self.operand.isConstExpression()
}

type CompiledBinaryExpression struct {
	CompiledBaseExpression
	left     CompiledExpression
	operator token.Token
	right    CompiledExpression
}

func (self CompiledBinaryExpression) isConstExpression() bool {
	return self.left.isConstExpression() && self.right.isConstExpression()
}

type CompiledAssignExpression struct {
	CompiledBaseExpression
	left     CompiledExpression
	operator token.Token
	right    CompiledExpression
}

func (self CompiledAssignExpression) isConstExpression() bool {
	return false
}

type CompiledFunLiteralExpression struct {
	CompiledBaseExpression
	funDefinition   string
	name            *ast.Identifier
	parameterList   *ast.ParameterList
	body            ast.Statement
	declarationList []*ast.VariableDeclaration
}

func (self CompiledFunLiteralExpression) isConstExpression() bool {
	return false
}

func (self *Compiler) createCompiledBaseExpression(index file.Index) CompiledBaseExpression {
	return CompiledBaseExpression{self, int(index) - 1}
}

func (self *Compiler) compileExpression(expression ast.Expression) CompiledExpression {
	switch expr := expression.(type) {
	case nil:
		return nil
	case *ast.NullLiteral:
		return self.compileNullLiteral(expr)
	case *ast.NumberLiteral:
		return self.compileNumberLiteral(expr)
	case *ast.StringLiteral:
		return self.compileStringLiteral(expr)
	case *ast.Identifier:
		return self.compileIdentifier(expr)
	case *ast.UnaryExpression:
		return self.compileUnaryExpression(expr)
	case *ast.BinaryExpression:
		return self.compileBinaryExpression(expr)
	case *ast.AssignExpression:
		return self.compileAssignExpression(expr)
	case *ast.FunLiteral:
		return self.compileFunLiteral(expr)
	default:
		return self.throwSyntaxError(int(expression.StartIndex())-1, "Unknown expression type: %T", expression)
	}
}

func (self *Compiler) compileNullLiteral(expr *ast.NullLiteral) CompiledExpression {
	return &CompiledLiteralExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
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
		return self.throwSyntaxError(int(expr.StartIndex())-1, "Unsupported number literal type: %T", expr.Value)
	}
	return &CompiledLiteralExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		value,
	}
}

func (self *Compiler) compileStringLiteral(expr *ast.StringLiteral) CompiledExpression {
	return &CompiledLiteralExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		ToStringValue(expr.Value),
	}
}

func (self *Compiler) compileIdentifier(expr *ast.Identifier) CompiledExpression {
	return &CompiledIdentifierExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		expr.Name,
	}
}

func (self *Compiler) compileUnaryExpression(expr *ast.UnaryExpression) CompiledExpression {
	return &CompiledUnaryExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		expr.Operator,
		self.compileExpression(expr.Operand),
		expr.Postfix,
	}
}

func (self *Compiler) compileBinaryExpression(expr *ast.BinaryExpression) CompiledExpression {
	return &CompiledBinaryExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		self.compileExpression(expr.Left),
		expr.Operator,
		self.compileExpression(expr.Right),
	}
}

func (self *Compiler) compileAssignExpression(expr *ast.AssignExpression) CompiledExpression {
	return &CompiledAssignExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		self.compileExpression(expr.Left),
		expr.Operator,
		self.compileExpression(expr.Right),
	}
}

func (self *Compiler) compileFunLiteral(expr *ast.FunLiteral) CompiledExpression {
	return &CompiledFunLiteralExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		expr.FunDefinition,
		expr.Name,
		expr.ParameterList,
		expr.Body,
		expr.DeclarationList,
	}
}
