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
			var s = 1
		}
		println(s)
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
