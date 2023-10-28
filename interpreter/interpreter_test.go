package interpreter

import (
	"fmt"
	"os"
	"testing"
)

func TestInterpreter(t *testing.T) {
	interpreter := CreateInterpreter()
	content, _ := os.ReadFile("../example/example.dl")
	value := interpreter.run("example.dl", string(content))
	println(fmt.Sprintf("%v", value.getVal()))
}
