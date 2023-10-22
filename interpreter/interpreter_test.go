package interpreter

import (
	"testing"
)

func TestInterpreter(t *testing.T) {
	//content, _ := os.ReadFile("../example/example.dl")
	interpreter := CreateInterpreter()
	value := interpreter.run("", `
		if 1 > 1 {
			return 1
		} else if 100 > 100 {
			return 2
		} else {
			return 3
		}
		return 4
	`)
	println(value.getValue())
}
