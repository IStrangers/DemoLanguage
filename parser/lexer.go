package parser

import (
	"github.com/istrangers/demolanguage/file"
	"github.com/istrangers/demolanguage/token"
	"strings"
)

func (parser *Parser) scan() (tkn token.Token, literal string, value string, index file.Index) {
	for {
		skipWhiteSpace := parser.skipWhiteSpace
		if skipWhiteSpace {
			parser.skipWhiteSpaceChr()
		}
		index = parser.IndexOf(parser.chrOffset)
		switch chr := parser.chr; {
		case !skipWhiteSpace && isWhiteSpaceChr(chr):
			tkn, literal, value = token.WHITE_SPACE, string(chr), string(chr)
			parser.readChr()
			break
		case isIdentifierStart(chr):
			literal = parser.scanIdentifier()
			value = literal
			keywordToken, exists := token.IsKeyword(literal)
			if exists {
				tkn = keywordToken
			} else {
				tkn = token.IDENTIFIER
			}
			break
		case isStringSymbol(chr):
			parser.readChr()
			value = parser.scanString()
			literal = string(chr) + value + string(parser.chr)
			tkn = token.STRING
			parser.readChr()
			break
		case isNumeric(chr):
			literal = parser.scanNumericLiteral()
			value = literal
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
				literal = tkn.String()
				value = tkn.String()
				break
			case '-':
				tkn = parser.switchToken(">,-,=", token.ARROW, token.DECREMENT, token.SUBTRACT_ASSIGN, token.SUBTRACT)
				literal = tkn.String()
				value = tkn.String()
				break
			case '*':
				tkn = parser.switchToken("=", token.MULTIPLY_ASSIGN, token.MULTIPLY)
				literal = tkn.String()
				value = tkn.String()
				break
			case '/':
				tkn = parser.switchToken("/,*,=", token.COMMENT, token.MULTI_COMMENT, token.DIVIDE_ASSIGN, token.DIVIDE)
				literal = tkn.String()
				value = tkn.String()
				if tkn == token.COMMENT || tkn == token.MULTI_COMMENT {
					comment := parser.scanComment(tkn)
					value = comment
					if tkn == token.COMMENT {
						comment = "//" + comment
					} else {
						comment = "/*" + comment + "*/"
						parser.readChr()
					}
					literal = comment
					if parser.skipComment {
						continue
					}
				}
				break
			case '%':
				tkn = parser.switchToken("=", token.REMAINDER_ASSIGN, token.REMAINDER)
				literal = tkn.String()
				value = tkn.String()
				break
			case '(':
				tkn, literal, value = token.LEFT_PARENTHESIS, string(chr), string(chr)
				break
			case ')':
				tkn, literal, value = token.RIGHT_PARENTHESIS, string(chr), string(chr)
				break
			case '{':
				tkn, literal, value = token.LEFT_BRACE, string(chr), string(chr)
				break
			case '}':
				tkn, literal, value = token.RIGHT_BRACE, string(chr), string(chr)
				break
			case '[':
				tkn, literal, value = token.LEFT_BRACKET, string(chr), string(chr)
				break
			case ']':
				tkn, literal, value = token.RIGHT_BRACKET, string(chr), string(chr)
				break
			case '.':
				tkn, literal, value = token.DOT, string(chr), string(chr)
				break
			case ',':
				tkn, literal, value = token.COMMA, string(chr), string(chr)
				break
			case ':':
				tkn, literal, value = token.COLON, string(chr), string(chr)
				break
			case ';':
				tkn, literal, value = token.SEMICOLON, string(chr), string(chr)
				break
			case '!':
				tkn = parser.switchToken("=", token.NOT_EQUAL, token.NOT)
				literal = tkn.String()
				value = tkn.String()
				break
			case '=':
				tkn = parser.switchToken("=", token.EQUAL, token.ASSIGN)
				literal = tkn.String()
				value = tkn.String()
				break
			case '<':
				tkn = parser.switchToken("=", token.LESS_OR_EQUAL, token.LESS)
				literal = tkn.String()
				value = tkn.String()
				break
			case '>':
				tkn = parser.switchToken("=", token.GREATER_OR_EQUAL, token.GREATER)
				literal = tkn.String()
				value = tkn.String()
				break
			case '&':
				tkn = parser.switchToken("&,=", token.LOGICAL_AND, token.AND_ARITHMETIC_ASSIGN, token.AND_ARITHMETIC)
				literal = tkn.String()
				value = tkn.String()
				break
			case '|':
				tkn = parser.switchToken("|,=", token.LOGICAL_OR, token.OR_ARITHMETIC_ASSIGN, token.OR_ARITHMETIC)
				literal = tkn.String()
				value = tkn.String()
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

func (parser *Parser) skipWhiteSpaceChr() {
	for isWhiteSpaceChr(parser.chr) {
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
			return multiCommentCount > 0 && chr != -1
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

func isWhiteSpaceChr(chr rune) bool {
	return chr == ' ' || chr == '\t' || chr == '\r' || chr == '\n' || chr == '\f'
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
	case '\u000a', '\u000d', '\u2028', '\u2029', -1:
		return true
	}
	return false
}
func isNotLineTerminator(chr rune) bool {
	return !isLineTerminator(chr)
}
