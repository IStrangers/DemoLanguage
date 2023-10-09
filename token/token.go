package token

type Token int

const (
	_ Token = iota
	ILLEGAL
	EOF

	IDENTIFIER
	ASSIGN

	NUMBER

	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE

	VAR
	FUN
)

var keywordMap = map[string]Token{
	"var": VAR,
	"fun": FUN,
}

func IsKeyword(k string) (Token, bool) {
	tkn, exists := keywordMap[k]
	return tkn, exists
}
