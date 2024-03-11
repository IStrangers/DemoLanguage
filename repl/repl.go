package main

import (
	"bufio"
	"fmt"
	"github.com/istrangers/demolanguage/vm"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	vm := vm.CreateVM()
	for {
		fmt.Print(">")
		content, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		result, err := vm.RunScript(content)
		if err != nil {
			panic(err)
		}
		if result == nil {
			continue
		}
		fmt.Printf("%v\n", result)
	}
}
