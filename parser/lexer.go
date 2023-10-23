package parser

import (
	"DemoLanguage/file"
	"DemoLanguage/token"
	"strings"
)

func (parser *Parser) scan() (tkn token.Token, literal string, index file.Index) {
	for {
		parser.skipWhiteSpace()
		index = parser.IndexOf(parser.chrOffset)
		switch chr := parser.chr; {
		case isIdentifierStart(chr):
			literal = parser.scanIdentifier()
			keywordToken, exists := token.IsKeyword(literal)
			if exists {
				tkn = keywordToken
			} else {
				tkn = token.IDENTIFIER
			}
			break
		case isStringSymbol(chr):
			parser.readChr()
			literal = parser.scanString()
			tkn = token.STRING
			parser.readChr()
			break
		case isNumeric(chr):
			literal = parser.scanNumericLiteral()
			tkn = token.NUMBER
			break
		default:
			parser.readChr()
			switch chr {
			case -1:
				tkn = token.EOF
				break
			case '+':
				tkn = parser.switchToken("+,=", token.INCREMENT, token.ADDITION_ASSIGN, token.ADDITION)
				break
			case '-':
				tkn = parser.switchToken(">,-,=", token.ARROW, token.DECREMENT, token.SUBTRACT_ASSIGN, token.SUBTRACT)
				break
			case '*':
				tkn = parser.switchToken("=", token.MULTIPLY_ASSIGN, token.MULTIPLY)
				break
			case '/':
				tkn = parser.switchToken("/,*,=", token.COMMENT, token.MULTI_COMMENT, token.DIVIDE_ASSIGN, token.DIVIDE)
				if tkn == token.COMMENT || tkn == token.MULTI_COMMENT {
					comment := parser.scanComment(tkn)
					if parser.skipComment {
						parser.readChr()
						continue
					}
					literal = comment
				}
				break
			case '%':
				tkn = parser.switchToken("=", token.REMAINDER_ASSIGN, token.REMAINDER)
				break
			case '(':
				tkn = token.LEFT_PARENTHESIS
				break
			case ')':
				tkn = token.RIGHT_PARENTHESIS
				break
			case '{':
				tkn = token.LEFT_BRACE
				break
			case '}':
				tkn = token.RIGHT_BRACE
				break
			case '[':
				tkn = token.LEFT_BRACKET
				break
			case ']':
				tkn = token.RIGHT_BRACKET
				break
			case '.':
				tkn = token.DOT
				break
			case ',':
				tkn = token.COMMA
				break
			case ':':
				tkn = token.COLON
				break
			case ';':
				tkn = token.SEMICOLON
				break
			case '!':
				tkn = parser.switchToken("=", token.NOT_EQUAL, token.NOT)
				break
			case '=':
				tkn = parser.switchToken("=", token.EQUAL, token.ASSIGN)
				break
			case '<':
				tkn = parser.switchToken("=", token.LESS_OR_EQUAL, token.LESS)
				break
			case '>':
				tkn = parser.switchToken("=", token.GREATER_OR_EQUEAL, token.GREATER)
				break
			case '&':
				tkn = parser.switchToken("&", token.LOGICAL_AND, token.AND_ARITHMETIC)
				break
			case '|':
				tkn = parser.switchToken("|", token.LOGICAL_OR, token.OR_ARITHMETIC)
				break
			default:
				tkn = token.ILLEGAL
				parser.errorUnexpected(index, tkn)
				break
			}
		}
		return
	}
}

func (parser *Parser) skipWhiteSpace() {
	for parser.chr == ' ' || parser.chr == '\t' || parser.chr == '\r' || parser.chr == '\n' || parser.chr == '\f' {
		parser.readChr()
	}
}

func (parser *Parser) readChr() rune {
	if parser.offset < parser.length {
		parser.chrOffset = parser.offset
		parser.chr = rune(parser.content[parser.offset])
		parser.offset += 1
		return parser.chr
	}
	parser.chrOffset = parser.length
	parser.chr = -1
	return parser.chr
}

func (parser *Parser) scanByFilter(filter func(rune) bool) string {
	chrOffset := parser.chrOffset
	for filter(parser.chr) {
		parser.readChr()
	}
	return parser.content[chrOffset:parser.chrOffset]
}

func (parser *Parser) scanIdentifier() string {
	return parser.scanByFilter(isIdentifierPart)
}

func (parser *Parser) scanNumericLiteral() string {
	return parser.scanByFilter(isNumericPart)
}

func (parser *Parser) scanString() string {
	return parser.scanByFilter(isNotStringSymbol)
}

func (parser *Parser) scanComment(tkn token.Token) string {
	if tkn == token.MULTI_COMMENT {
		multiCommentCount := 1
		multiComment := parser.scanByFilter(func(chr rune) bool {
			if chr == '/' && parser.readChr() == '*' {
				multiCommentCount++
			}
			if chr == '*' && parser.readChr() == '/' {
				multiCommentCount--
			}
			return multiCommentCount > 0
		})
		return multiComment[:len(multiComment)-1]
	} else {
		return parser.scanByFilter(isNotLineTerminator)
	}
}

func (parser *Parser) switchToken(keysStr string, tkns ...token.Token) token.Token {
	keys := strings.Split(keysStr, ",")
	for i, key := range keys {
		if parser.chr == rune(key[0]) {
			parser.readChr()
			return tkns[i]
		}
	}
	return tkns[len(tkns)-1]
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
func isNumericPart(chr rune) bool {
	return chr == '.' || isNumeric(chr)
}

func isStringSymbol(chr rune) bool {
	return chr == '"' || chr == '\''
}
func isNotStringSymbol(chr rune) bool {
	return !isStringSymbol(chr)
}

func isLineTerminator(chr rune) bool {
	switch chr {
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	}
	return false
}
func isNotLineTerminator(chr rune) bool {
	return !isLineTerminator(chr)
}
