package parser

import (
	"DemoLanguage/ast"
	"DemoLanguage/file"
	"DemoLanguage/token"
)

func (parser *Parser) parseStatementList() []ast.Statement {
	return parser.parseStatementListByCondition(func(tkn token.Token) bool {
		return tkn != token.EOF
	})
}

func (parser *Parser) parseStatementListByCondition(endCondition func(token.Token) bool) []ast.Statement {
	var statementList []ast.Statement
	for endCondition(parser.token) {
		statementList = append(statementList, parser.parseStatement())
	}
	return statementList
}

func (parser *Parser) nextStatement() {

}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.token {
	case token.EOF:
		parser.errorUnexpectedToken(parser.token)
		return &ast.BadStatement{Start: parser.index, End: parser.index + 1}
	case token.VAR:
		return parser.parseVarStatement()
	case token.FUN:
		return parser.parseFunStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	case token.IF:
		return parser.parseIfStatement()
	case token.FOR:
		return parser.parseForStatement()
	case token.SWITCH:
		return parser.parseSwitchStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseBlockStatement() ast.Statement {
	return &ast.BlockStatement{
		LeftBrace: parser.expect(token.LEFT_BRACE),
		Body: parser.parseStatementListByCondition(func(tkn token.Token) bool {
			return tkn != token.RIGHT_BRACE && tkn != token.EOF
		}),
		RightBrace: parser.expect(token.RIGHT_BRACE),
	}
}

func (parser *Parser) parseVarStatement() ast.Statement {
	varIndex := parser.expect(token.VAR)
	list := parser.parseVarDeclarationList(varIndex)
	return &ast.VarStatement{
		Var:  varIndex,
		List: list,
	}
}

func (parser *Parser) parseVarDeclarationList(varIndex file.Index) []*ast.Binding {
	bindingList := parser.parseBindingList()

	parser.scope.AddDeclaration(&ast.VariableDeclaration{
		Var:  varIndex,
		List: bindingList,
	})

	return bindingList
}

func (parser *Parser) parseFunStatement() ast.Statement {
	funStatement := &ast.FunStatement{
		Fun:        parser.expect(token.FUN),
		Name:       parser.parseIdentifier(),
		FunLiteral: parser.parseFunLiteral(),
	}
	funStatement.FunDefinition = parser.slice(funStatement.StartIndex(), funStatement.FunLiteral.EndIndex())
	return funStatement
}

func (parser *Parser) parseReturnStatement() ast.Statement {
	return &ast.ReturnStatement{
		Return:    parser.expect(token.RETURN),
		Arguments: parser.parseReturnArguments(),
	}
}

func (parser *Parser) parseIfStatement() ast.Statement {
	ifStatement := &ast.IfStatement{
		If:        parser.expect(token.IF),
		Condition: parser.parseExpression(),
	}
	ifStatement.Consequent = parser.parseBlockStatement()
	if parser.token == token.ELSE {
		parser.next()
		ifStatement.Alternate = parser.parseStatement()
	}
	return ifStatement
}

func (parser *Parser) parseForStatement() ast.Statement {
	return nil
}

func (parser *Parser) parseSwitchStatement() ast.Statement {
	switchStatement := &ast.SwitchStatement{
		Switch:       parser.expect(token.SWITCH),
		Discriminant: parser.parseExpression(),
		Default:      -1,
	}
	switchStatement.Body, switchStatement.Default = parser.parseCaseStatementList()
	switchStatement.RightBrace = parser.expect(token.RIGHT_BRACE)
	return switchStatement
}

func (parser *Parser) parseCaseStatementList() ([]*ast.CaseStatement, int) {
	return nil, 0
}

func (parser *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Expression: nil,
	}
}
