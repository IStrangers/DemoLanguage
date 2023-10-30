package vm

import "DemoLanguage/file"

type Program struct {
	values       ValueArray
	instructions InstructionArray
	file         *file.File
}
