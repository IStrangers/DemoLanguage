package parser

import "DemoLanguage/token"

type Parser struct {
	lexer Lexer
}

func (parser *Parser) Parse() {
	lexer := parser.lexer
	for lexer.token != token.EOF {
		tkn, literal := lexer.scan()
		println(tkn.String(), literal)
	}
}

func CreateParser(content string) Parser {
	return Parser{
		lexer: Lexer{
			content:   content,
			length:    len(content),
			chr:       rune(content[0]),
			chrOffset: 0,
		},
	}
}
