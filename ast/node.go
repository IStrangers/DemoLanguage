package ast

import (
	"DemoLanguage/file"
)

type Program struct {
	Body []Statement
	File *file.File
}

type Node interface {
	StartIndex() file.Index
	EndIndex() file.Index
}

type (
	Statement interface {
		Node
	}

	BadStatement struct {
		Start file.Index
		End   file.Index
	}

	ExpressionStatement struct {
		Expression Expression
	}

	BlockStatement struct {
		LeftBrace  file.Index
		Body       []Statement
		RightBrace file.Index
	}

	VarStatement struct {
		Var  file.Index
		List []*Binding
	}

	FunStatement struct {
		FunDefinition string
		Fun           file.Index
		Name          *Identifier
		FunLiteral    *FunLiteral
	}

	ReturnStatement struct {
		Return    file.Index
		Arguments []Expression
	}

	IfStatement struct {
		If         file.Index
		Condition  Expression
		Consequent Statement
		Alternate  Statement
	}
)

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
	return self.RightBrace
}

func (self *VarStatement) StartIndex() file.Index {
	return self.Var
}
func (self *VarStatement) EndIndex() file.Index {
	return self.List[len(self.List)-1].EndIndex()
}

func (self *FunStatement) StartIndex() file.Index {
	return self.Fun
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

type (
	Expression interface {
		Node
	}

	BindingTarget interface {
		Expression
	}

	Binding struct {
		Target      BindingTarget
		Initializer Expression
	}

	Identifier struct {
		Index file.Index
		Name  string
	}

	BadExpression struct {
		Start file.Index
		End   file.Index
	}

	NumberLiteral struct {
		Index   file.Index
		Literal string
		Value   any
	}

	StringLiteral struct {
		Index   file.Index
		Literal string
		Value   string
	}

	ParameterList struct {
		LeftParenthesis  file.Index
		List             []*Binding
		RightParenthesis file.Index
	}

	FunLiteral struct {
		FunDefinition   string
		ParameterList   *ParameterList
		Body            Statement
		DeclarationList []*VariableDeclaration
	}
)

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

func (self *BadExpression) StartIndex() file.Index {
	return self.Start
}
func (self *BadExpression) EndIndex() file.Index {
	return self.End
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

func (self *ParameterList) StartIndex() file.Index {
	return self.LeftParenthesis
}
func (self *ParameterList) EndIndex() file.Index {
	return self.RightParenthesis
}

func (self *FunLiteral) StartIndex() file.Index {
	return self.ParameterList.StartIndex()
}
func (self *FunLiteral) EndIndex() file.Index {
	return self.Body.EndIndex() + 1
}

type (
	VariableDeclaration struct {
		Var  file.Index
		List []*Binding
	}
)
