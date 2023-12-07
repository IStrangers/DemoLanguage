package vm

import (
	"DemoLanguage/parser"
	"fmt"
	"testing"
)

func TestCompiler(t *testing.T) {
	parser := parser.CreateParser(1, "", `
		//var a = 600
		//if a > 1000 {
		//	a = 1
		//} else if(a > 500) {
		//	a = 2
		//} else {
		//	a = 3
		//}
		//switch a {
		//	case 1 {
		//		a = 4
		//	}
		//	case 2 {
		//		a = 5
		//	}
		//	default {
		//		a = 6
		//	}
		//}
		//var i = 1
		//for ;i <= 5; {
		//	i = 6
		//}
		//a(100)
		//fun a(b,c = 1) {
		//	var a = 1
		//	return a + b + c
		//}
		//fun getFebNum(n) {
		//	if n <= 2 {
		//		return 1
		//	} else {
		//		return getFebNum(n - 1) + getFebNum(n - 2)
		//	}
		//}
		//getFebNum(20)
		//var a = 2
		//if a == 2 && a > 1 {
		//	a = 1
		//} else {
		//	a = 2
		//}
		var obj = {
			a: 1,
			b: "123",
			c: {
				a: 2
			},
			d: fun(a) {
				return  a + 100
			}
		}
		var arr = [1, obj.c.a, obj["c"]['a'],obj.d]
		obj.d(66)
		arr[3](99)
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
			maxCallStackSize: 999,
		},
	}
	compiler.compile(program)
	evalVM := compiler.evalVM
	evalVM.program = compiler.program
	evalVM.runTry()
	result := evalVM.result
	fmt.Printf("%v\n", result)
}
