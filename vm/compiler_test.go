package vm

import (
	"DemoLanguage/parser"
	"fmt"
	"testing"
)

func TestCompiler(t *testing.T) {
	parser := parser.CreateParser(1, "", `
		var a = 600
		if a > 1000 {
			a = 1
		} else if a > 500 {
			a = 2
		} else {
			a = 3
		}
		switch a {
			case 1 {
				a = 4
			}
			case 2 {
				a = 5
			}
			default {
				a = 6
			}
		}
		var i = 1
		for ;i <= 5; {
			i = 6
		}
		b(100)
		fun b(b,c = 1) {
			var a = 1
			return a + b + c
		}
		fun getFebNum(n) {
			if n <= 2 {
				return 1
			} else {
				return getFebNum(n - 1) + getFebNum(n - 2)
			}
		}
		getFebNum(20)
		var a = 2
		if a == 2 && a > 1 {
			a = 1
		} else {
			a = 2
		}
		var obj = {
			a: 1,
			b: "123",
			c: {
				a: 2
			},
			d: fun(a) {
				return this.a + 100 + a
			},
			e: (p1) -> {
				return p1 * 2 * this.c.a
			},
			f: p1 -> p1 / 2
		}
		var arr = [1, obj.c.a, obj["c"]['a'],obj.d]
		println(obj)
		println(arr)
		obj.d(10)
		obj["d"](20)
		arr[3](30)
		obj.e(999)
		obj.f(1000)
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
