package vm

import (
	"fmt"
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	content, _ := os.ReadFile("../example/vm_example.dl")
	vm := CreateVM()
	result, err := vm.RunScript(string(content))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", result)
}
