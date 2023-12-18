package parser

import (
	"fmt"
	"github.com/istrangers/demolanguage/file"
	"github.com/istrangers/demolanguage/token"
)

const (
	ERR_UnexpectedToken      = "Unexpected token %v"
	ERR_UnexpectedEndOfInput = "Unexpected end of input"
)

type Error struct {
	Message  string
	Position *file.Position
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

func (errorList *ErrorList) Errors() error {
	if errorList.Length() > 0 {
		return errorList
	}
	return nil
}

func (errorList *ErrorList) Error() string {
	length := errorList.Length()
	switch length {
	case 0:
		return "no errors"
	case 1:
		return errorList.FirstError().Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", errorList.FirstError().Error(), length)
}

func (errorList *ErrorList) AddError(error *Error) {
	*errorList = append(*errorList, error)
}
func (errorList *ErrorList) Add(message string, position *file.Position) {
	errorList.AddError(&Error{Message: message, Position: position})
}

func (errorList ErrorList) Length() int {
	return len(errorList)
}

func (errorList ErrorList) FirstError() *Error {
	if errorList.Length() > 0 {
		return errorList[0]
	}
	return nil
}

func (errorList ErrorList) LastError() *Error {
	return errorList[errorList.Length()-1]
}

func (parser *Parser) errorUnexpected(index file.Index, tkn token.Token) error {
	return parser.error(index, ERR_UnexpectedToken, tkn)
}

func (parser *Parser) errorUnexpectedToken(tkn token.Token) error {
	message := ERR_UnexpectedToken
	messageValue := tkn.String()
	switch tkn {
	case token.EOF:
		message = ERR_UnexpectedEndOfInput
		break
	}
	return parser.error(parser.index, message, messageValue)
}

func (parser *Parser) error(index file.Index, message string, messageValues ...any) *Error {
	position := parser.Position(index)
	message = fmt.Sprintf(message, messageValues...)
	parser.errors.Add(message, position)
	return parser.errors.LastError()
}
