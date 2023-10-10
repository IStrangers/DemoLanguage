package parser

import "DemoLanguage/token"

type Parser struct {
	baseOffset int
	file       *File

	content   string
	length    int
	chr       rune
	chrOffset int
	offset    int
	token     token.Token
	literal   string

	errors ErrorList
}

func CreateParser(baseOffset int, fileName string, content string) *Parser {
	return &Parser{
		baseOffset: baseOffset,
		file:       CreateFile(baseOffset, fileName, content),
		content:    content,
		length:     len(content),
		chr:        ' ',
	}
}

func (parser *Parser) Parse() {

}
