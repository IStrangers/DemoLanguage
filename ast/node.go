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

	BlockStatement struct {
		Statement
		LeftBrace  file.Index
		Body       []Statement
		RightBrace file.Index
	}
)

type (
	Expression interface {
		Node
	}
)

func (self *BlockStatement) StartIndex() file.Index {
	return self.LeftBrace
}
func (self *BlockStatement) EndIndex() file.Index {
	return self.RightBrace + 1
}
