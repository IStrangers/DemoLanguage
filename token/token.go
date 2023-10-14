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
	IDENTIFIER

	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )
	LEFT_BRACE        // {
	RIGHT_BRACE       // }
	LEFT_BRACKET      // [
	RIGHT_BRACKET     // ]
	PERIOD            // .
	COMMA             // ,
	COLON             // :
	SEMICOLON         // ;

	NUMBER
	STRING

	ADDITION         // +
	SUBTRACT         // -
	MULTIPLY         // *
	DIVIDE           // /
	REMAINDER        // %
	AND_ARITHMETIC   // &
	OR_ARITHMETIC    // |
	INCREMENT        // ++
	DECREMENT        // --
	ADDITION_ASSIGN  // +=
	SUBTRACT_ASSIGN  // -=
	MULTIPLY_ASSIGN  // *=
	DIVIDE_ASSIGN    // /=
	REMAINDER_ASSIGN // %=

	ASSIGN            // =
	EQUAL             // ==
	NOT               // ！
	NOT_EQUAL         // !=
	LESS              // <
	LESS_OR_EQUAL     // <=
	GREATER           // >
	GREATER_OR_EQUEAL // >=
	LOGICAL_AND       // &&
	LOGICAL_OR        // ||

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
)

var tokenStringMap = [...]string{
	0:          "UNKNOWN",
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
	IDENTIFIER: "IDENTIFIER",

	LEFT_PARENTHESIS:  "(",
	RIGHT_PARENTHESIS: ")",
	LEFT_BRACE:        "{",
	RIGHT_BRACE:       "}",
	LEFT_BRACKET:      "[",
	RIGHT_BRACKET:     "]",
	PERIOD:            ".",
	COMMA:             ",",
	COLON:             ":",
	SEMICOLON:         ";",

	NUMBER: "NUMBER",
	STRING: "STRING",

	ADDITION:         "+",
	SUBTRACT:         "-",
	MULTIPLY:         "*",
	DIVIDE:           "/",
	REMAINDER:        "%",
	AND_ARITHMETIC:   "&",
	OR_ARITHMETIC:    "|",
	INCREMENT:        "++",
	DECREMENT:        "--",
	ADDITION_ASSIGN:  "+=",
	SUBTRACT_ASSIGN:  "-=",
	MULTIPLY_ASSIGN:  "*=",
	DIVIDE_ASSIGN:    "/=",
	REMAINDER_ASSIGN: "%=",

	ASSIGN:            "=",
	EQUAL:             "==",
	NOT:               "!",
	NOT_EQUAL:         "!=",
	LESS:              "<",
	LESS_OR_EQUAL:     "<=",
	GREATER:           ">",
	GREATER_OR_EQUEAL: ">=",
	LOGICAL_AND:       "&&",
	LOGICAL_OR:        "||",

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
}

var keywordMap = map[string]Token{
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
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}
