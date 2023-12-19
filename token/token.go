package token

type Token int

func (tkn Token) String() string {
	return tokenStringMap[tkn]
}

const (
	_ Token = iota
	ILLEGAL
	EOF
	COMMENT
	MULTI_COMMENT
	WHITE_SPACE
	IDENTIFIER

	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )
	LEFT_BRACE        // {
	RIGHT_BRACE       // }
	LEFT_BRACKET      // [
	RIGHT_BRACKET     // ]
	DOT               // .
	COMMA             // ,
	COLON             // :
	SEMICOLON         // ;
	ARROW             // ->

	NUMBER
	STRING
	BOOLEAN
	NULL

	ADDITION              // +
	SUBTRACT              // -
	MULTIPLY              // *
	DIVIDE                // /
	REMAINDER             // %
	AND_ARITHMETIC        // &
	OR_ARITHMETIC         // |
	INCREMENT             // ++
	DECREMENT             // --
	ADDITION_ASSIGN       // +=
	SUBTRACT_ASSIGN       // -=
	MULTIPLY_ASSIGN       // *=
	DIVIDE_ASSIGN         // /=
	REMAINDER_ASSIGN      // %=
	AND_ARITHMETIC_ASSIGN // &=
	OR_ARITHMETIC_ASSIGN  // |=

	ASSIGN           // =
	EQUAL            // ==
	NOT              // ！
	NOT_EQUAL        // !=
	LESS             // <
	LESS_OR_EQUAL    // <=
	GREATER          // >
	GREATER_OR_EQUAL // >=
	LOGICAL_AND      // &&
	LOGICAL_OR       // ||

	VAR    // var
	FUN    // fun
	RETURN // return

	IF       // if
	ELSE     // else
	BREAK    // break
	FOR      // for
	SWITCH   // switch
	CASE     // case
	DEFAULT  // default
	CONTINUE // continue
	THIS     // this
)

var tokenStringMap = [...]string{
	0:             "UNKNOWN",
	ILLEGAL:       "ILLEGAL",
	EOF:           "EOF",
	COMMENT:       "COMMENT",
	MULTI_COMMENT: "MULTI_COMMENT",
	WHITE_SPACE:   "WHITE_SPACE",
	IDENTIFIER:    "IDENTIFIER",

	LEFT_PARENTHESIS:  "(",
	RIGHT_PARENTHESIS: ")",
	LEFT_BRACE:        "{",
	RIGHT_BRACE:       "}",
	LEFT_BRACKET:      "[",
	RIGHT_BRACKET:     "]",
	DOT:               ".",
	COMMA:             ",",
	COLON:             ":",
	SEMICOLON:         ";",
	ARROW:             "->",

	NUMBER:  "NUMBER",
	STRING:  "STRING",
	BOOLEAN: "BOOLEAN",
	NULL:    "NULL",

	ADDITION:              "+",
	SUBTRACT:              "-",
	MULTIPLY:              "*",
	DIVIDE:                "/",
	REMAINDER:             "%",
	AND_ARITHMETIC:        "&",
	OR_ARITHMETIC:         "|",
	INCREMENT:             "++",
	DECREMENT:             "--",
	ADDITION_ASSIGN:       "+=",
	SUBTRACT_ASSIGN:       "-=",
	MULTIPLY_ASSIGN:       "*=",
	DIVIDE_ASSIGN:         "/=",
	REMAINDER_ASSIGN:      "%=",
	AND_ARITHMETIC_ASSIGN: "&=",
	OR_ARITHMETIC_ASSIGN:  "|=",

	ASSIGN:           "=",
	EQUAL:            "==",
	NOT:              "!",
	NOT_EQUAL:        "!=",
	LESS:             "<",
	LESS_OR_EQUAL:    "<=",
	GREATER:          ">",
	GREATER_OR_EQUAL: ">=",
	LOGICAL_AND:      "&&",
	LOGICAL_OR:       "||",

	VAR:    "var",
	FUN:    "fun",
	RETURN: "return",

	IF:       "if",
	ELSE:     "else",
	BREAK:    "break",
	FOR:      "for",
	SWITCH:   "switch",
	CASE:     "case",
	DEFAULT:  "default",
	CONTINUE: "continue",
	THIS:     "this",
}

var keywordMap = map[string]Token{
	"true":     BOOLEAN,
	"false":    BOOLEAN,
	"null":     NULL,
	"var":      VAR,
	"fun":      FUN,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"break":    BREAK,
	"for":      FOR,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"continue": CONTINUE,
	"this":     THIS,
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}
