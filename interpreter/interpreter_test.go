package interpreter

import (
	"fmt"
	"testing"
)

func TestInterpreter(t *testing.T) {
	//content, _ := os.ReadFile("../example/example.dl")
	interpreter := CreateInterpreter()
	value := interpreter.run("", `
		/*var a = 500
		if a == 100 {
			return 1
		} else if a == 200 {
			return 2
		}
		switch a {
			case 300 {
				return 3
			}
			case 400 {
				return 4
			}
			/*default {
				return 6
			}*/
		}
		for var index = 1;index < 100;index++ {
			if index % 5 == 0 || (index >= 90 && index != 95) {
				if index == 5 {
					continue
				}
				a = index
				break
			}
		}
		return a*/
		var c = 99.99
		fun add(a = 1,b = 1) {
			return a + b + c
		}
		add(50,100)
	`)
	println(fmt.Sprintf("%v", value.getVal()))
}
