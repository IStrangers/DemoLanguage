package parser

import (
	"DemoLanguage/token"
	"fmt"
)

const (
	ERR_UnexpectedToken      = "Unexpected token %v"
	ERR_UnexpectedEndOfInput = "Unexpected end of input"
)

type Error struct {
	Message  string
	Position *Position
}

func (error *Error) Error() string {
	position := error.Position
	fileName := position.FileName
	if fileName == "" {
		fileName = "(anonymous)"
	}
	return fmt.Sprintf(
		"%s Line %d:%d %s",
		fileName,
		position.Line,
		position.Column,
		error.Message,
	)
}

type ErrorList []*Error

func (errorList *ErrorList) AddError(error *Error) {
	*errorList = append(*errorList, error)
}
func (errorList *ErrorList) Add(message string, position *Position) {
	errorList.AddError(&Error{Message: message, Position: position})
}

func (errorList ErrorList) Len() int {
	return len(errorList)
}

func (errorList ErrorList) LastError() *Error {
	return errorList[errorList.Len()-1]
}

func (parser *Parser) errorUnexpected(index Index, tkn token.Token) error {
	return parser.error(index, ERR_UnexpectedToken, tkn)
}

func (parser *Parser) error(index Index, message string, tkn token.Token) *Error {
	position := parser.Position(index)
	message = fmt.Sprintf(message, tkn)
	parser.errors.Add(message, position)
	return parser.errors.LastError()
}
