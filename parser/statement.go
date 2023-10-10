package parser

import (
	"DemoLanguage/ast"
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

func (parser *Parser) parseStatement() ast.Statement {
	return nil
}

func (parser *Parser) parseBlockStatementList() *ast.BlockStatement {
	return &ast.BlockStatement{
		LeftBrace: parser.expect(token.LEFT_BRACE),
		Body: parser.parseStatementListByCondition(func(tkn token.Token) bool {
			return tkn != token.RIGHT_BRACE && tkn != token.EOF
		}),
		RightBrace: parser.expect(token.RIGHT_BRACE),
	}
}
