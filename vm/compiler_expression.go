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

type CompiledObjectLiteralExpression struct {
	CompiledBaseExpression
	properties []ast.Property
}

func (self CompiledObjectLiteralExpression) isConstExpression() bool {
	return false
}

type CompiledArrayLiteralExpression struct {
	CompiledBaseExpression
	values []ast.Expression
}

func (self CompiledArrayLiteralExpression) isConstExpression() bool {
	return false
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
	operator := self.operator
	if operator == token.LOGICAL_OR || operator == token.LOGICAL_AND {
		if !self.left.isConstExpression() {
			return false
		}
		if v, ex := self.compile.evalConstExpr(self.left); ex == nil {
			boolVal := v.toBool()
			if (operator == token.LOGICAL_OR && boolVal) || (operator == token.LOGICAL_AND && !boolVal) {
				return true
			}
			return self.right.isConstExpression()
		}
		return true
	}
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

type CompiledCallExpression struct {
	CompiledBaseExpression
	callee    CompiledExpression
	arguments []CompiledExpression
}

func (self CompiledCallExpression) isConstExpression() bool {
	return false
}

type CompiledDotExpression struct {
	CompiledBaseExpression
	left CompiledExpression
	name string
}

func (self CompiledDotExpression) isConstExpression() bool {
	return false
}

type CompiledBracketExpression struct {
	CompiledBaseExpression
	left        CompiledExpression
	indexOrName CompiledExpression
}

func (self CompiledBracketExpression) isConstExpression() bool {
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
	case *ast.ObjectLiteral:
		return self.compileObjectLiteral(expr)
	case *ast.ArrayLiteral:
		return self.compileArrayLiteral(expr)
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
	case *ast.CallExpression:
		return self.compileCallExpression(expr)
	case *ast.DotExpression:
		return self.compileDotExpression(expr)
	case *ast.BracketExpression:
		return self.compileBracketExpression(expr)
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

func (self *Compiler) compileObjectLiteral(expr *ast.ObjectLiteral) CompiledExpression {
	return &CompiledObjectLiteralExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		expr.Properties,
	}
}

func (self *Compiler) compileArrayLiteral(expr *ast.ArrayLiteral) CompiledExpression {
	return &CompiledArrayLiteralExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		expr.Values,
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

func (self *Compiler) compileCallExpression(expr *ast.CallExpression) CompiledExpression {
	var arguments []CompiledExpression
	for _, argument := range expr.Arguments {
		arguments = append(arguments, self.compileExpression(argument))
	}
	return &CompiledCallExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		self.compileExpression(expr.Callee),
		arguments,
	}
}

func (self *Compiler) compileDotExpression(expr *ast.DotExpression) CompiledExpression {
	return &CompiledDotExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		self.compileExpression(expr.Left),
		expr.Identifier.Name,
	}
}

func (self *Compiler) compileBracketExpression(expr *ast.BracketExpression) CompiledExpression {
	return &CompiledBracketExpression{
		self.createCompiledBaseExpression(expr.StartIndex()),
		self.compileExpression(expr.Left),
		self.compileExpression(expr.Expression),
	}
}
