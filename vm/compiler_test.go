package vm

import (
	"DemoLanguage/parser"
	"fmt"
	"testing"
)

func TestCompiler(t *testing.T) {
	//content, _ := os.ReadFile("../example/example.dl")
	parser := parser.CreateParser(1, "example.dl", `
		fun a() {
			var df = "东风"
			return fun() {
				return fun() {
					return df
				}
			}	
		}
		println(a()()())
	`)
	program, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}
	compiler := CreateCompiler()
	compiler.compile(program)
	evalVM := compiler.evalVM
	evalVM.program = compiler.program
	evalVM.program.dumpInstructions(t.Logf)
	evalVM.runTry()
	result := evalVM.result
	fmt.Printf("%v\n", result)
}
