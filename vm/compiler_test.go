package vm

import (
	"DemoLanguage/parser"
	"testing"
)

func TestCompiler(t *testing.T) {
	parser := parser.CreateParser(1, "", `
		var a = 100
		if a > 50 {
			a = 1
		} else {
			a = 2
		}
	`)
	program, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}
	compiler := &Compiler{
		program: &Program{},
		evalVM: &VM{
			runtime: &Runtime{
				globalObject: &Object{self: &BaseObject{
					valueMapping: make(map[string]Value),
				}},
			},
		},
	}
	compiler.compile(program)
	compiler.evalVM.program = compiler.program
	compiler.evalVM.runTry()
}
