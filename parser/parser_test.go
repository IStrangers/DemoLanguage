package parser

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	content, _ := os.ReadFile("../example/parser_example.dl")
	parser := CreateParser(1, "", string(content))
	program, err := parser.Parse()
	if err != nil {
		println(err.Error())
		return
	}
	println(program)
}
