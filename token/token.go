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

	ASSIGN            // =
	EQUAL             // ==
	NOT               // ！
	NOT_EQUAL         // !=
	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )
	LEFT_BRACE        // {
	RIGHT_BRACE       // }

	NUMBER
	STRING

	ADD       // +
	SUBTRACT  // -
	MULTIPLY  // *
	DIVIDE    // /
	REMAINDER // %

	VAR    // var
	FUN    // fun
	RETURN // return

	IF     // if
	FOR    // for
	SWITCH // switch
)

var tokenStringMap = [...]string{
	0:          "UNKNOWN",
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
	IDENTIFIER: "IDENTIFIER",

	ASSIGN:            "=",
	EQUAL:             "==",
	NOT:               "!",
	NOT_EQUAL:         "!=",
	LEFT_PARENTHESIS:  "(",
	RIGHT_PARENTHESIS: ")",
	LEFT_BRACE:        "{",
	RIGHT_BRACE:       "}",

	NUMBER: "NUMBER",
	STRING: "STRING",

	ADD:       "+",
	SUBTRACT:  "-",
	MULTIPLY:  "*",
	DIVIDE:    "/",
	REMAINDER: "%",

	VAR:    "var",
	FUN:    "fun",
	RETURN: "return",

	IF:     "if",
	FOR:    "for",
	SWITCH: "switch",
}

var keywordMap = map[string]Token{
	"var":    VAR,
	"fun":    FUN,
	"return": RETURN,
	"if":     IF,
	"for":    FOR,
	"switch": SWITCH,
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}
