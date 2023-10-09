package parser

import (
	"DemoLanguage/token"
)

type Lexer struct {
	content   string
	length    int
	chr       rune
	chrOffset int
	token     token.Token
	literal   string
}

func (lexer *Lexer) scan() (tkn token.Token, literal string) {
	for {
		lexer.skipWhiteSpace()
		switch chr := lexer.chr; {
		case isIdentifierStart(chr):
			literal := lexer.scanIdentifier()
			keywordToken, exists := token.IsKeyword(literal)
			if exists {
				tkn = keywordToken
			} else {
				tkn = token.IDENTIFIER
			}
			break
		case isNumeric(chr):
			literal = lexer.scanNumericLiteral()
			tkn = token.NUMBER
			break
		default:
			switch chr {
			case '+':
				tkn = token.ADD
				break
			case '-':
				tkn = token.SUBTRACT
				break
			case '*':
				tkn = token.MULTIPLY
				break
			case '/':
				tkn = token.DIVIDE
				break
			case '=':
				tkn = token.ASSIGN
				lexer.readChr()
				break
			}
		}
		return
	}
}

func (lexer *Lexer) skipWhiteSpace() {
	for lexer.chr == ' ' || lexer.chr == '\t' || lexer.chr == '\r' || lexer.chr == '\n' || lexer.chr == '\f' {
		lexer.readChr()
	}
}

func (lexer *Lexer) readChr() {
	pos := lexer.chrOffset + 1
	if pos < lexer.length {
		lexer.chr = rune(lexer.content[pos])
		lexer.chrOffset = pos
		return
	}
	lexer.token = token.EOF
	lexer.chrOffset = lexer.length
	lexer.chr = -1
}

func (lexer *Lexer) scanByFilter(filter func(rune) bool) string {
	chrOffset := lexer.chrOffset
	for filter(lexer.chr) {
		lexer.readChr()
	}
	return lexer.content[chrOffset:lexer.chrOffset]
}

func (lexer *Lexer) scanIdentifier() string {
	return lexer.scanByFilter(isIdentifierPart)
}

func (lexer *Lexer) scanNumericLiteral() string {
	return lexer.scanByFilter(isNumeric)
}

func isIdentifierStart(chr rune) bool {
	return chr == '$' || chr == '_' || (chr >= 'A' && chr <= 'Z') || (chr >= 'a' && chr <= 'z')
}
func isIdentifierPart(chr rune) bool {
	return isIdentifierStart(chr) || isNumeric(chr)
}
func isNumeric(chr rune) bool {
	return chr >= '0' && chr <= '9'
}
