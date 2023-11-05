package vm

import (
	"DemoLanguage/parser"
	"testing"
)

func TestCompiler(t *testing.T) {
	parser := parser.CreateParser(1, "", `
		var a = 100
		if a > 500 {
			a = 1
		} else if(a > 1000) {
			a = 2
		} else {
			a = 3
		}
		switch 1 {
			case 1 {
				a = 4
			}
			case 2 {
				a = 5
			}
		}
	`)
	program, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}
	compiler := &Compiler{
		program: &Program{
			source: program.File,
		},
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
	result := compiler.evalVM.result
	if result == nil {
		return
	}
	println(result.toString())
}
