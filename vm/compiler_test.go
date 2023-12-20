package vm

import (
	"fmt"
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	content, _ := os.ReadFile("../example/example.dl")
	vm := CreateVM()
	result, _ := vm.RunScript(string(content))
	fmt.Printf("%v\n", result)
}
