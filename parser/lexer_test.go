package parser

import (
	"DemoLanguage/token"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	content, _ := os.ReadFile("../example/lexer_example.dl")
	parser := CreateParser(1, "", string(content))
	for parser.token != token.EOF {
		tkn, literal := parser.scan()
		println(tkn.String(), literal)
	}
	for _, err := range parser.errors {
		println(err.Error())
	}
}
