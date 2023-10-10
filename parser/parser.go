package parser

import (
	"DemoLanguage/ast"
	"DemoLanguage/file"
	"DemoLanguage/token"
)

type Parser struct {
	baseOffset int
	file       *file.File

	content   string
	length    int
	chr       rune
	chrOffset int
	offset    int
	token     token.Token
	literal   string
	index     file.Index

	errors ErrorList
}

func CreateParser(baseOffset int, fileName string, content string) *Parser {
	return &Parser{
		baseOffset: baseOffset,
		file:       file.CreateFile(baseOffset, fileName, content),
		content:    content,
		length:     len(content),
		chr:        ' ',
	}
}

func (parser *Parser) ParseProgram() *ast.Program {
	return &ast.Program{
		Body: parser.parseStatementList(),
		File: parser.file,
	}
}

func (parser *Parser) Parse() (*ast.Program, error) {
	program := parser.ParseProgram()
	return program, &parser.errors
}

func (parser *Parser) next() {
	parser.token, parser.literal, parser.index = parser.scan()
}

func (parser *Parser) expect(tkn token.Token) file.Index {
	if parser.token != tkn {
		parser.errorUnexpectedToken(tkn)
	}
	parser.next()
	return parser.index
}

func (parser *Parser) IndexOf(offset int) file.Index {
	return file.Index(parser.baseOffset + offset)
}

func (parser *Parser) Position(index file.Index) *file.Position {
	return parser.file.Position(int(index) - parser.baseOffset)
}
