package vm

import (
	"DemoLanguage/parser"
	"testing"
)

func TestCompiler(t *testing.T) {
	parser := parser.CreateParser(1, "", `
		if 50 > 1 {
		}
	`)
	program, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}
	compiler := &Compiler{
		program: &Program{},
		evalVM:  &VM{},
	}
	compiler.compile(program)
}
