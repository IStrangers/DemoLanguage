package main

import "DemoLanguage/parser"

func main() {
	content := `
		var a = 1
	`
	parser := parser.CreateParser(content)
	parser.Parse()
}
