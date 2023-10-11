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

	scope *Scope
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
	parser.next()
	return &ast.Program{
		Body: parser.parseStatementList(),
		File: parser.file,
	}
}

func (parser *Parser) Parse() (*ast.Program, error) {
	parser.openScope()
	defer parser.closeScope()
	program := parser.ParseProgram()
	return program, parser.errors.Errors()
}

func (parser *Parser) next() {
	parser.token, parser.literal, parser.index = parser.scan()
}

func (parser *Parser) expect(tkn token.Token) file.Index {
	index := parser.index
	if parser.token != tkn {
		parser.errorUnexpectedToken(tkn)
	}
	parser.next()
	return index
}

func (parser *Parser) IndexOf(offset int) file.Index {
	return file.Index(parser.baseOffset + offset)
}

func (parser *Parser) Position(index file.Index) *file.Position {
	return parser.file.Position(int(index) - parser.baseOffset)
}

func (parser *Parser) slice(start, end file.Index) string {
	from := int(start) - parser.baseOffset
	to := int(end) - parser.baseOffset
	if from >= 0 && to <= len(parser.content) {
		return parser.content[from:to]
	}
	return ""
}
