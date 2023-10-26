package interpreter

import (
	"fmt"
	"testing"
)

func TestInterpreter(t *testing.T) {
	//content, _ := os.ReadFile("../example/example.dl")
	interpreter := CreateInterpreter()
	value := interpreter.run("", `
		var a = 500
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
		var c = 99.99
		fun add(a = 1,b = 1) {
			return a + b + c
		}
		println(add(100))
		fun getFebNum(n) {
			if n == 1 || n == 2 {
				return 1
			} else {
				return getFebNum(n - 1) + getFebNum(n - 2)
			}
		}
		println(getFebNum(20))
		var obj = {
			name: "Afghanistan",
			count: add(),
			ref: c,
			inline: {
				ref: a
			},
			arrowFun1: (p1) -> {
				println(p1 + " * 2")
				return p1 * 2
			},
			arrowFun2: p1 -> p1 / 2
		}
		var arr = [obj.inline.ref,"sfd",500]
		arr.add(555)
		arr.remove(2)
		var v = arr.get(arr.size() - 1)
		arr[0] = v
		obj["name"] = 123
		obj.name = "564654"
		println(obj)
		println(arr)
		var arrowFun = p1 -> p1 / 2
		println(arrowFun(5))
		return obj.arrowFun1(500)
	`)
	println(fmt.Sprintf("%v", value.getVal()))
}
