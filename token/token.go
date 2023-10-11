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
	COMMA             // ,

	NUMBER
	STRING

	ADD            // +
	SUBTRACT       // -
	MULTIPLY       // *
	DIVIDE         // /
	REMAINDER      // %
	AND_ARITHMETIC // &
	OR_ARITHMETIC  // |

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

	IF     // if
	ELSE   // else
	FOR    // for
	SWITCH // switch
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
	COMMA:             ",",

	NUMBER: "NUMBER",
	STRING: "STRING",

	ADD:            "+",
	SUBTRACT:       "-",
	MULTIPLY:       "*",
	DIVIDE:         "/",
	REMAINDER:      "%",
	AND_ARITHMETIC: "&",
	OR_ARITHMETIC:  "|",

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

	IF:     "if",
	ELSE:   "else",
	FOR:    "for",
	SWITCH: "switch",
}

var keywordMap = map[string]Token{
	"var":    VAR,
	"fun":    FUN,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"switch": SWITCH,
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}
