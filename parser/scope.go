package parser

import "DemoLanguage/ast"

type Scope struct {
	outer           *Scope
	declarationList []*ast.VariableDeclaration
	labels          []string

	inSwitch    bool
	inIteration bool
	inFunction  bool
}

func (scpoe *Scope) AddDeclaration(declaration *ast.VariableDeclaration) {
	scpoe.declarationList = append(scpoe.declarationList, declaration)
}

func (parser *Parser) openScope() {
	parser.scope = &Scope{
		outer: parser.scope,
	}
	scope := parser.scope
	outer := parser.scope.outer
	if outer != nil {
		scope.inSwitch, scope.inIteration, scope.inFunction = outer.inSwitch, outer.inIteration, outer.inFunction
	}
}

func (parser *Parser) closeScope() {
	parser.scope = parser.scope.outer
}
