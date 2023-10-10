package main

import (
	"DemoLanguage/parser"
	"os"
)

func main() {
	content, _ := os.ReadFile("./example/example1.dl")
	parser := parser.CreateParser(string(content))
	parser.Parse()
}
