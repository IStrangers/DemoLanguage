package parser

import (
	"DemoLanguage/token"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	content, _ := os.ReadFile("../example/example.dl")
	parser := CreateParser(1, "", string(content))
	for parser.token != token.EOF {
		parser.next()
		println(parser.token.String(), parser.literal, parser.index)
	}
	for _, err := range parser.errors {
		println(err.Error())
	}
}
