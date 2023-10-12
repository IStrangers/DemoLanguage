package parser

import (
	"DemoLanguage/ast"
	"DemoLanguage/token"
)

func (parser *Parser) parseBindingList() (bindingList []*ast.Binding) {
	for {
		bindingList = append(bindingList, parser.parseBinding())
		if parser.token != token.COMMA {
			break
		}
		parser.next()
	}
	return
}

func (parser *Parser) parseBinding() *ast.Binding {
	binding := &ast.Binding{
		Target: parser.parseBindingTarget(),
	}

	if parser.token == token.ASSIGN {
		parser.next()
		binding.Initializer = parser.parseAssignExpression()
	}
	return binding
}

func (parser *Parser) parseBindingTarget() ast.BindingTarget {
	switch parser.token {
	case token.IDENTIFIER:
		return parser.parseIdentifier()
	default:
		index := parser.expect(token.IDENTIFIER)
		badExpression := &ast.BadExpression{
			Start: index,
			End:   index,
		}
		parser.nextStatement()
		return badExpression
	}
}

func (parser *Parser) parseIdentifier() *ast.Identifier {
	defer parser.expect(token.IDENTIFIER)
	return &ast.Identifier{
		Index: parser.index,
		Name:  parser.literal,
	}
}

func (parser *Parser) parseFunLiteral() *ast.FunLiteral {
	funLiteral := &ast.FunLiteral{}
	funLiteral.ParameterList = parser.parseFunParameterList()
	funLiteral.Body, funLiteral.DeclarationList = parser.parseFunBlock()
	funLiteral.FunDefinition = parser.slice(funLiteral.StartIndex(), funLiteral.EndIndex())
	return funLiteral
}

func (parser *Parser) parseFunParameterList() *ast.ParameterList {
	return &ast.ParameterList{
		LeftParenthesis:  parser.expect(token.LEFT_PARENTHESIS),
		List:             parser.parseBindingList(),
		RightParenthesis: parser.expect(token.RIGHT_PARENTHESIS),
	}
}

func (parser *Parser) parseFunBlock() (ast.Statement, []*ast.VariableDeclaration) {
	parser.openScope()
	defer parser.closeScope()
	return parser.parseBlockStatement(), parser.scope.declarationList
}

func (parser *Parser) parseReturnArguments() (arguments []ast.Expression) {
	for parser.token != token.RIGHT_BRACE {
		arguments = append(arguments, parser.parseExpression())
	}
	return
}

func (parser *Parser) parseExpression() ast.Expression {
	left := parser.parseAssignExpression()

	return left
}

func (parser *Parser) parseAssignExpression() ast.Expression {
	parenthesis := false

	switch parser.token {
	case token.LEFT_PARENTHESIS:
		parenthesis = true
	}

	left := parser.parseConditionalExpression()

	var operator token.Token
	switch parser.token {
	case token.ASSIGN:
		operator = token.ASSIGN
	}

	if operator != 0 {
		index := parser.index
		err := true

		switch left.(type) {
		case *ast.Identifier:
			err = false
			break
		case *ast.ArrayLiteral:
			if parenthesis || operator != token.ASSIGN {
				break
			}
			err = false
			break
		case *ast.ObjectLiteral:
			if parenthesis || operator != token.ASSIGN {
				break
			}
			err = false
			break
		}
		if err {
			parser.error(left.StartIndex(), "Invalid left-hand side in assignment")
			parser.nextStatement()
			return &ast.BadExpression{Start: index, End: parser.index}
		}
		return &ast.AssignExpression{
			Left:     left,
			Operator: operator,
			Right:    parser.parseAssignExpression(),
		}
	}

	return left
}

func (parser *Parser) parseConditionalExpression() ast.Expression {
	return nil
}

func (parser *Parser) parsePrimaryExpression() ast.Expression {
	index := parser.index

	switch parser.token {
	case token.IDENTIFIER:
		return parser.parseIdentifier()
	case token.NUMBER:
		return parser.parseNumberLiteral()
	case token.STRING:
		return parser.parseStringLiteral()
	}

	parser.errorUnexpectedToken(parser.token)
	return &ast.BadExpression{
		Start: index,
		End:   parser.index,
	}
}

func (parser *Parser) parseNumberLiteral() ast.Expression {
	defer parser.expect(token.NUMBER)
	return &ast.NumberLiteral{
		Index:   parser.index,
		Literal: parser.literal,
		Value:   parser.literal,
	}
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	defer parser.expect(token.STRING)
	return &ast.StringLiteral{
		Index:   parser.index,
		Literal: parser.literal,
		Value:   parser.literal,
	}
}
