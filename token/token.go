package token

type Token int

func (tkn Token) String() string {
	return tokenStringMap[tkn]
}

const (
	_ Token = iota
	ILLEGAL
	EOF

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

	ADD      // +
	SUBTRACT // -
	MULTIPLY // *
	DIVIDE   // /

	VAR    // var
	FUN    // fun
	RETURN // return
)

var tokenStringMap = [...]string{
	0:       "UNKNOWN",
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

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

	ADD:      "+",
	SUBTRACT: "-",
	MULTIPLY: "*",
	DIVIDE:   "/",

	VAR:    "var",
	FUN:    "fun",
	RETURN: "return",
}

var keywordMap = map[string]Token{
	"var":    VAR,
	"fun":    FUN,
	"return": RETURN,
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}
