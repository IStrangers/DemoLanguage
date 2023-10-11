package parser

import "DemoLanguage/ast"

type Scope struct {
	outer           *Scope
	declarationList []*ast.VariableDeclaration
	labels          []string
}

func (scpoe *Scope) AddDeclaration(declaration *ast.VariableDeclaration) {
	scpoe.declarationList = append(scpoe.declarationList, declaration)
}

func (parser *Parser) openScope() {
	parser.scope = &Scope{
		outer: parser.scope,
	}
}

func (parser *Parser) closeScope() {
	parser.scope = parser.scope.outer
}
