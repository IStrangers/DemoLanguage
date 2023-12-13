package vm

import (
	"DemoLanguage/parser"
	"fmt"
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	content, _ := os.ReadFile("../example/example.dl")
	parser := parser.CreateParser(1, "example.dl", string(content))
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
