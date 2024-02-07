package ast

import (
	"github.com/istrangers/demolanguage/file"
	"github.com/istrangers/demolanguage/token"
)

type Program struct {
	Body            []Statement
	DeclarationList []*VariableDeclaration
	File            *file.File
}

type Node interface {
	StartIndex() file.Index
	EndIndex() file.Index
}

type (
	Statement interface {
		Node
		statement()
	}

	AbstractStatement struct {
		Statement
	}

	BadStatement struct {
		AbstractStatement
		Start file.Index
		End   file.Index
	}

	ExpressionStatement struct {
		AbstractStatement
		Expression Expression
	}

	BlockStatement struct {
		AbstractStatement
		LeftBrace  file.Index
		Body       []Statement
		RightBrace file.Index
	}

	VarStatement struct {
		AbstractStatement
		Var  file.Index
		List []*Binding
	}

	FunStatement struct {
		AbstractStatement
		FunLiteral *FunLiteral
	}

	ReturnStatement struct {
		AbstractStatement
		Return    file.Index
		Arguments []Expression
	}

	IfStatement struct {
		AbstractStatement
		If         file.Index
		Condition  Expression
		Consequent Statement
		Alternate  Statement
	}

	ForStatement struct {
		AbstractStatement
		For         file.Index
		Initializer Statement
		Condition   Expression
		Update      Expression
		Body        Statement
	}

	SwitchStatement struct {
		AbstractStatement
		Switch       file.Index
		Discriminant Expression
		Body         []*CaseStatement
		Default      int
		RightBrace   file.Index
	}

	CaseStatement struct {
		AbstractStatement
		Case       file.Index
		Condition  Expression
		Consequent Statement
	}

	BranchStatement interface {
		Statement
		Token() token.Token
	}

	BreakStatement struct {
		AbstractStatement
		Break file.Index
	}

	ContinueStatement struct {
		AbstractStatement
		Continue file.Index
	}

	ThrowStatement struct {
		AbstractStatement
		Throw    file.Index
		Argument Expression
	}

	TryCatchFinallyStatement struct {
		AbstractStatement
		Try             file.Index
		TryBody         *BlockStatement
		CatchParameters *ParameterList
		CatchBody       *BlockStatement
		FinallyBody     *BlockStatement
	}
)

func (self *AbstractStatement) statement() {
}

func (self *BadStatement) StartIndex() file.Index {
	return self.Start
}

func (self *BadStatement) EndIndex() file.Index {
	return self.End
}

func (self *ExpressionStatement) StartIndex() file.Index {
	return self.Expression.StartIndex()
}
func (self *ExpressionStatement) EndIndex() file.Index {
	return self.Expression.EndIndex()
}

func (self *BlockStatement) StartIndex() file.Index {
	return self.LeftBrace
}
func (self *BlockStatement) EndIndex() file.Index {
	return self.RightBrace + 1
}

func (self *VarStatement) StartIndex() file.Index {
	return self.Var
}
func (self *VarStatement) EndIndex() file.Index {
	return self.List[len(self.List)-1].EndIndex()
}

func (self *FunStatement) StartIndex() file.Index {
	return self.FunLiteral.StartIndex()
}
func (self *FunStatement) EndIndex() file.Index {
	return self.FunLiteral.EndIndex()
}

func (self *ReturnStatement) StartIndex() file.Index {
	return self.Return
}
func (self *ReturnStatement) EndIndex() file.Index {
	argsLength := len(self.Arguments)
	if argsLength > 0 {
		return self.Arguments[argsLength-1].EndIndex()
	}
	return self.Return + 6
}

func (self *IfStatement) StartIndex() file.Index {
	return self.If
}
func (self *IfStatement) EndIndex() file.Index {
	if self.Alternate != nil {
		return self.Alternate.EndIndex()
	}
	return self.Consequent.EndIndex()
}

func (self *ForStatement) StartIndex() file.Index {
	return self.For
}
func (self *ForStatement) EndIndex() file.Index {
	return self.Body.EndIndex()
}

func (self *SwitchStatement) StartIndex() file.Index {
	return self.Switch
}
func (self *SwitchStatement) EndIndex() file.Index {
	return self.RightBrace + 1
}

func (self *CaseStatement) StartIndex() file.Index {
	return self.Case
}
func (self *CaseStatement) EndIndex() file.Index {
	return self.Consequent.EndIndex()
}

func (self *BreakStatement) StartIndex() file.Index {
	return self.Break
}
func (self *BreakStatement) EndIndex() file.Index {
	return self.Break + 5
}
func (self *BreakStatement) Token() token.Token {
	return token.BREAK
}

func (self *ContinueStatement) StartIndex() file.Index {
	return self.Continue
}
func (self *ContinueStatement) EndIndex() file.Index {
	return self.Continue + 8
}
func (self *ContinueStatement) Token() token.Token {
	return token.CONTINUE
}

func (self *ThrowStatement) StartIndex() file.Index {
	return self.Throw
}
func (self *ThrowStatement) EndIndex() file.Index {
	if self.Argument != nil {
		return self.Argument.EndIndex()
	}
	return self.Throw + 5
}

func (self *TryCatchFinallyStatement) StartIndex() file.Index {
	return self.Try
}
func (self *TryCatchFinallyStatement) EndIndex() file.Index {
	if self.FinallyBody != nil {
		return self.FinallyBody.EndIndex()
	}
	return self.CatchBody.EndIndex()
}

type (
	Expression interface {
		Node
		expression()
	}

	AbstractExpression struct {
		Expression
	}

	Property interface {
		Expression
	}

	BindingTarget interface {
		Expression
	}

	Binding struct {
		AbstractExpression
		Target      BindingTarget
		Initializer Expression
	}

	Identifier struct {
		AbstractExpression
		Index file.Index
		Name  string
	}

	AssignExpression struct {
		AbstractExpression
		Left     Expression
		Operator token.Token
		Right    Expression
	}

	NumberLiteral struct {
		AbstractExpression
		Index   file.Index
		Literal string
		Value   any
	}

	StringLiteral struct {
		AbstractExpression
		Index   file.Index
		Literal string
		Value   string
	}

	BooleanLiteral struct {
		AbstractExpression
		Index file.Index
		Value bool
	}

	NullLiteral struct {
		AbstractExpression
		Index file.Index
	}

	ArrayLiteral struct {
		AbstractExpression
		LeftBracket  file.Index
		Values       []Expression
		RightBracket file.Index
	}

	ObjectLiteral struct {
		AbstractExpression
		LeftBrace  file.Index
		Properties []Property
		RightBrace file.Index
	}

	PropertyKeyValue struct {
		AbstractExpression
		Name  *Identifier
		Colon file.Index
		Value Expression
	}

	ParameterList struct {
		AbstractExpression
		LeftParenthesis  file.Index
		List             []*Binding
		RightParenthesis file.Index
	}

	FunLiteral struct {
		AbstractExpression
		Fun             file.Index
		Name            *Identifier
		ParameterList   *ParameterList
		Body            *BlockStatement
		DeclarationList []*VariableDeclaration
		FunDefinition   string
	}

	ArrowFunctionLiteral struct {
		AbstractExpression
		Index           file.Index
		ParameterList   *ParameterList
		Arrow           file.Index
		Body            *BlockStatement
		DeclarationList []*VariableDeclaration
		FunDefinition   string
	}

	BinaryExpression struct {
		AbstractExpression
		Operator   token.Token
		Left       Expression
		Right      Expression
		Comparison bool
	}

	UnaryExpression struct {
		AbstractExpression
		Index    file.Index
		Operator token.Token
		Operand  Expression
		Postfix  bool
	}

	ThisExpression struct {
		AbstractExpression
		Index file.Index
	}

	DotExpression struct {
		AbstractExpression
		Left       Expression
		Dot        file.Index
		Identifier *Identifier
	}

	BracketExpression struct {
		AbstractExpression
		Left         Expression
		LeftBracket  file.Index
		Expression   Expression
		RightBracket file.Index
	}

	CallExpression struct {
		AbstractExpression
		Callee           Expression
		LeftParenthesis  file.Index
		Arguments        []Expression
		RightParenthesis file.Index
	}

	NewExpression struct {
		AbstractExpression
		New              file.Index
		Callee           Expression
		LeftParenthesis  file.Index
		Arguments        []Expression
		RightParenthesis file.Index
	}

	BadExpression struct {
		AbstractExpression
		Start file.Index
		End   file.Index
	}
)

func (self *AbstractExpression) expression() {
}

func (self *Binding) StartIndex() file.Index {
	return self.Target.StartIndex()
}
func (self *Binding) EndIndex() file.Index {
	if self.Initializer != nil {
		return self.Initializer.EndIndex()
	}
	return self.Target.EndIndex()
}

func (self *Identifier) StartIndex() file.Index {
	return self.Index
}
func (self *Identifier) EndIndex() file.Index {
	return file.Index(int(self.Index) + len(self.Name))
}

func (self *AssignExpression) StartIndex() file.Index {
	return self.Left.StartIndex()
}
func (self *AssignExpression) EndIndex() file.Index {
	return self.Right.EndIndex()
}

func (self *NumberLiteral) StartIndex() file.Index {
	return self.Index
}
func (self *NumberLiteral) EndIndex() file.Index {
	return file.Index(int(self.Index) + len(self.Literal))
}

func (self *StringLiteral) StartIndex() file.Index {
	return self.Index
}
func (self *StringLiteral) EndIndex() file.Index {
	return file.Index(int(self.Index) + len(self.Literal))
}

func (self *BooleanLiteral) StartIndex() file.Index {
	return self.Index
}
func (self *BooleanLiteral) EndIndex() file.Index {
	if self.Value {
		return self.Index + 4
	}
	return self.Index + 5
}

func (self *NullLiteral) StartIndex() file.Index {
	return self.Index
}
func (self *NullLiteral) EndIndex() file.Index {
	return self.Index + 4
}

func (self *ArrayLiteral) StartIndex() file.Index {
	return self.LeftBracket
}
func (self *ArrayLiteral) EndIndex() file.Index {
	return self.RightBracket + 1
}

func (self *ObjectLiteral) StartIndex() file.Index {
	return self.LeftBrace
}
func (self *ObjectLiteral) EndIndex() file.Index {
	return self.RightBrace + 1
}

func (self *PropertyKeyValue) StartIndex() file.Index {
	return self.Name.StartIndex()
}
func (self *PropertyKeyValue) EndIndex() file.Index {
	return self.Value.EndIndex()
}

func (self *ParameterList) StartIndex() file.Index {
	return self.LeftParenthesis
}
func (self *ParameterList) EndIndex() file.Index {
	return self.RightParenthesis + 1
}

func (self *FunLiteral) StartIndex() file.Index {
	return self.Fun
}
func (self *FunLiteral) EndIndex() file.Index {
	return self.Body.EndIndex()
}

func (self *ArrowFunctionLiteral) StartIndex() file.Index {
	return self.Index
}
func (self *ArrowFunctionLiteral) EndIndex() file.Index {
	return self.Body.EndIndex()
}

func (self *BinaryExpression) StartIndex() file.Index {
	return self.Left.StartIndex()
}
func (self *BinaryExpression) EndIndex() file.Index {
	return self.Right.EndIndex()
}

func (self *UnaryExpression) StartIndex() file.Index {
	if self.Postfix {
		return self.Operand.StartIndex()
	}
	return self.Index
}
func (self *UnaryExpression) EndIndex() file.Index {
	if self.Postfix {
		return self.Operand.EndIndex() + 2
	}
	return self.Operand.EndIndex()
}

func (self *ThisExpression) StartIndex() file.Index {
	return self.Index
}
func (self *ThisExpression) EndIndex() file.Index {
	return self.Index + 4
}

func (self *DotExpression) StartIndex() file.Index {
	return self.Left.StartIndex()
}
func (self *DotExpression) EndIndex() file.Index {
	return self.Identifier.EndIndex()
}

func (self *BracketExpression) StartIndex() file.Index {
	return self.Left.StartIndex()
}
func (self *BracketExpression) EndIndex() file.Index {
	return self.RightBracket + 1
}

func (self *CallExpression) StartIndex() file.Index {
	return self.Callee.StartIndex()
}
func (self *CallExpression) EndIndex() file.Index {
	return self.RightParenthesis + 1
}

func (self *NewExpression) StartIndex() file.Index {
	return self.New
}
func (self *NewExpression) EndIndex() file.Index {
	return self.Callee.EndIndex()
}

func (self *BadExpression) StartIndex() file.Index {
	return self.Start
}
func (self *BadExpression) EndIndex() file.Index {
	return self.End
}

type AccessModifier int

const (
	_ AccessModifier = iota
	PRIVATE
	PROTECTED
	PUBLIC
)

type (
	VariableDeclaration struct {
		Var  file.Index
		List []*Binding
	}

	Declaration interface {
		declaration()
	}

	AbstractDeclaration struct {
		Declaration
	}

	InterfaceDeclaration struct {
		AbstractStatement
		AbstractDeclaration
		Index      file.Index
		LeftBrace  file.Index
		Body       []Declaration
		RightBrace file.Index
	}

	ClassDeclaration struct {
		AbstractStatement
		AbstractDeclaration
		Index      file.Index
		Name       *Identifier
		SuperClass *Identifier
		Interfaces []*Identifier
		LeftBrace  file.Index
		Body       []Declaration
		RightBrace file.Index
	}

	StaticBlockDeclaration struct {
		AbstractStatement
		AbstractDeclaration
		Static     file.Index
		LeftBrace  file.Index
		Body       *BlockStatement
		RightBrace file.Index
	}

	FieldDeclaration struct {
		AbstractStatement
		AbstractDeclaration
		Index          file.Index
		AccessModifier AccessModifier
		Static         bool
		Name           *Identifier
		Initializer    Expression
	}

	MethodDeclaration struct {
		AbstractStatement
		AbstractDeclaration
		Index          file.Index
		AccessModifier AccessModifier
		Static         bool
		Body           *FunLiteral
	}
)

func (self *AbstractDeclaration) declaration() {
}

func (self *InterfaceDeclaration) StartIndex() file.Index {
	return self.Index
}
func (self *InterfaceDeclaration) EndIndex() file.Index {
	return self.RightBrace
}

func (self *ClassDeclaration) StartIndex() file.Index {
	return self.Index
}
func (self *ClassDeclaration) EndIndex() file.Index {
	return self.RightBrace
}

func (self *StaticBlockDeclaration) StartIndex() file.Index {
	return self.Static
}
func (self *StaticBlockDeclaration) EndIndex() file.Index {
	return self.RightBrace
}

func (self *FieldDeclaration) StartIndex() file.Index {
	return self.Index
}
func (self *FieldDeclaration) EndIndex() file.Index {
	if self.Initializer != nil {
		return self.Initializer.EndIndex()
	}
	return self.Name.EndIndex()
}

func (self *MethodDeclaration) StartIndex() file.Index {
	return self.Index
}
func (self *MethodDeclaration) EndIndex() file.Index {
	return self.Body.EndIndex()
}
