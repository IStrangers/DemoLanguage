package parser

import (
	"github.com/istrangers/demolanguage/ast"
	"github.com/istrangers/demolanguage/file"
	"github.com/istrangers/demolanguage/token"
)

type Parser struct {
	baseOffset     int
	file           *file.File
	skipComment    bool
	skipWhiteSpace bool

	content   string
	length    int
	chr       rune
	chrOffset int
	offset    int
	token     token.Token
	literal   string
	value     string
	index     file.Index

	errors ErrorList

	scope *Scope
}

func CreateParser(baseOffset int, fileName string, content string, skipComment bool, skipWhiteSpace bool) *Parser {
	return &Parser{
		baseOffset:     baseOffset,
		file:           file.CreateFile(baseOffset, fileName, content),
		skipComment:    skipComment,
		skipWhiteSpace: skipWhiteSpace,
		content:        content,
		length:         len(content),
		chr:            ' ',
	}
}

func (parser *Parser) Parse() (*ast.Program, error) {
	parser.openScope()
	defer parser.closeScope()
	program := parser.parseProgram()
	return program, parser.errors.Errors()
}

func (parser *Parser) ScanNextToken() (token.Token, string, string, file.Index) {
	return parser.scan()
}

func (parser *Parser) parseProgram() *ast.Program {
	parser.next()
	return &ast.Program{
		Body:            parser.parseStatementList(),
		DeclarationList: parser.scope.declarationList,
		File:            parser.file,
	}
}

func (parser *Parser) next() {
	parser.token, parser.literal, parser.value, parser.index = parser.scan()
}

func (parser *Parser) expect(tkn token.Token) file.Index {
	index := parser.index
	if parser.token != tkn {
		parser.errorUnexpectedToken(parser.token)
	}
	parser.next()
	return index
}

func (parser *Parser) expectToken(tkn token.Token) token.Token {
	if parser.token != tkn {
		parser.errorUnexpectedToken(parser.token)
	}
	parser.next()
	return tkn
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

type ParseState struct {
	chr        rune
	chrOffset  int
	offset     int
	token      token.Token
	literal    string
	value      string
	index      file.Index
	errorIndex int
}

func (parser *Parser) markParseState() *ParseState {
	return &ParseState{
		chr:        parser.chr,
		chrOffset:  parser.chrOffset,
		offset:     parser.offset,
		token:      parser.token,
		literal:    parser.literal,
		value:      parser.value,
		index:      parser.index,
		errorIndex: parser.errors.Length(),
	}
}

func (parser *Parser) restoreParseState(parseState *ParseState) {
	parser.chr = parseState.chr
	parser.chrOffset = parseState.chrOffset
	parser.offset = parseState.offset
	parser.token = parseState.token
	parser.literal = parseState.literal
	parser.value = parseState.value
	parser.index = parseState.index
	parser.errors = parser.errors[:parseState.errorIndex]
}
