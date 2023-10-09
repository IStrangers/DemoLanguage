package parser

import "DemoLanguage/token"

type Parser struct {
	lexer Lexer
}

func (parser *Parser) Parse() {
	lexer := parser.lexer
	for lexer.token != token.EOF {
		println(lexer.scan())
	}
}

func CreateParser(content string) Parser {
	return Parser{
		lexer: Lexer{
			content:   content,
			length:    len(content),
			chr:       ' ',
			chrOffset: 0,
		},
	}
}
