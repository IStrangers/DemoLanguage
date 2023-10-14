package parser

import (
	"DemoLanguage/ast"
	"DemoLanguage/file"
	"DemoLanguage/token"
	"strconv"
)

func (parser *Parser) parseBindingList() (bindingList []*ast.Binding) {
	for {
		bindingList = append(bindingList, parser.parseBinding())
		if parser.token != token.COMMA {
			break
		}
		parser.next()
	}
	return
}

func (parser *Parser) parseBinding() *ast.Binding {
	binding := &ast.Binding{
		Target: parser.parseBindingTarget(),
	}

	if parser.token == token.ASSIGN {
		parser.next()
		binding.Initializer = parser.parseAssignExpression()
	}
	return binding
}

func (parser *Parser) parseBindingTarget() ast.BindingTarget {
	switch parser.token {
	case token.IDENTIFIER:
		return parser.parseIdentifier()
	default:
		index := parser.expect(token.IDENTIFIER)
		badExpression := &ast.BadExpression{
			Start: index,
			End:   index,
		}
		parser.nextStatement()
		return badExpression
	}
}

func (parser *Parser) parseIdentifier() *ast.Identifier {
	defer parser.expect(token.IDENTIFIER)
	return &ast.Identifier{
		Index: parser.index,
		Name:  parser.literal,
	}
}

func (parser *Parser) parseFunLiteral() *ast.FunLiteral {
	funLiteral := &ast.FunLiteral{}
	funLiteral.ParameterList = parser.parseFunParameterList()
	funLiteral.Body, funLiteral.DeclarationList = parser.parseFunBlock()
	funLiteral.FunDefinition = parser.slice(funLiteral.StartIndex(), funLiteral.EndIndex())
	return funLiteral
}

func (parser *Parser) parseFunParameterList() *ast.ParameterList {
	return &ast.ParameterList{
		LeftParenthesis:  parser.expect(token.LEFT_PARENTHESIS),
		List:             parser.parseBindingList(),
		RightParenthesis: parser.expect(token.RIGHT_PARENTHESIS),
	}
}

func (parser *Parser) parseFunBlock() (ast.Statement, []*ast.VariableDeclaration) {
	parser.openScope()
	defer parser.closeScope()
	return parser.parseBlockStatement(), parser.scope.declarationList
}

func (parser *Parser) parseReturnArguments() (arguments []ast.Expression) {
	for parser.token != token.RIGHT_BRACE {
		arguments = append(arguments, parser.parseExpression())
	}
	return
}

func (parser *Parser) parseExpression() ast.Expression {
	left := parser.parseAssignExpression()

	return left
}

func (parser *Parser) parseAssignExpression() ast.Expression {
	parenthesis := false

	switch parser.token {
	case token.LEFT_PARENTHESIS:
		parenthesis = true
	}

	left := parser.parseConditionalExpression()

	var operator token.Token
	switch parser.token {
	case token.ASSIGN:
		operator = token.ASSIGN
	}

	if operator != 0 {
		index := parser.index
		err := true

		switch left.(type) {
		case *ast.Identifier:
			err = false
			break
		case *ast.ArrayLiteral:
			if parenthesis || operator != token.ASSIGN {
				break
			}
			err = false
			break
		case *ast.ObjectLiteral:
			if parenthesis || operator != token.ASSIGN {
				break
			}
			err = false
			break
		}
		if err {
			parser.error(left.StartIndex(), "Invalid left-hand side in assignment")
			parser.nextStatement()
			return &ast.BadExpression{Start: index, End: parser.index}
		}
		parser.expect(operator)
		return &ast.AssignExpression{
			Left:     left,
			Operator: operator,
			Right:    parser.parseAssignExpression(),
		}
	}

	return left
}

func (parser *Parser) parseConditionalExpression() ast.Expression {
	left := parser.parseLogicalOrExpression()

	return left
}

func (parser *Parser) parseLogicalOrExpression() ast.Expression {
	left := parser.parseLogicalAndExpression()

	switch parser.token {
	case token.LOGICAL_OR:
		binaryExpression := &ast.BinaryExpression{
			Operator: parser.expectToken(parser.token),
			Left:     left,
			Right:    parser.parseLogicalAndExpression(),
		}
		return binaryExpression
	}
	return left
}

func (parser *Parser) parseLogicalAndExpression() ast.Expression {
	left := parser.parseBitwiseOrExpression()

	switch parser.token {
	case token.LOGICAL_AND:
		binaryExpression := &ast.BinaryExpression{
			Operator: parser.expectToken(parser.token),
			Left:     left,
			Right:    parser.parseBitwiseOrExpression(),
		}
		return binaryExpression
	}
	return left
}

func (parser *Parser) parseBitwiseOrExpression() ast.Expression {
	left := parser.parseBitwiseExclusiveOrExpression()

	return left
}

func (parser *Parser) parseBitwiseExclusiveOrExpression() ast.Expression {
	left := parser.parseBitwiseAndExpression()

	return left
}

func (parser *Parser) parseBitwiseAndExpression() ast.Expression {
	left := parser.parseEqualityExpression()

	return left
}

func (parser *Parser) parseEqualityExpression() ast.Expression {
	left := parser.parseRelationalExpression()

	switch parser.token {
	case token.EQUAL, token.NOT_EQUAL:
		binaryExpression := &ast.BinaryExpression{
			Operator:   parser.expectToken(parser.token),
			Left:       left,
			Right:      parser.parseRelationalExpression(),
			Comparison: true,
		}
		return binaryExpression
	}
	return left
}

func (parser *Parser) parseRelationalExpression() ast.Expression {
	left := parser.parseShiftExpression()

	switch parser.token {
	case token.LESS, token.LESS_OR_EQUAL, token.GREATER, token.GREATER_OR_EQUEAL:
		binaryExpression := &ast.BinaryExpression{
			Operator:   parser.expectToken(parser.token),
			Left:       left,
			Right:      parser.parseShiftExpression(),
			Comparison: true,
		}
		return binaryExpression
	}
	return left
}

func (parser *Parser) parseShiftExpression() ast.Expression {
	left := parser.parseAdditiveExpression()

	return left
}

func (parser *Parser) parseAdditiveExpression() ast.Expression {
	left := parser.parseMultiplicativeExpression()

	switch parser.token {
	case token.ADDITION, token.SUBTRACT:
		binaryExpression := &ast.BinaryExpression{
			Operator: parser.expectToken(parser.token),
			Left:     left,
			Right:    parser.parseMultiplicativeExpression(),
		}
		return binaryExpression
	}
	return left
}

func (parser *Parser) parseMultiplicativeExpression() ast.Expression {
	left := parser.parseExponentiationExpression()

	switch parser.token {
	case token.MULTIPLY, token.DIVIDE, token.REMAINDER:
		binaryExpression := &ast.BinaryExpression{
			Operator: parser.expectToken(parser.token),
			Left:     left,
			Right:    parser.parseExponentiationExpression(),
		}
		return binaryExpression
	}
	return left
}

func (parser *Parser) parseExponentiationExpression() ast.Expression {
	left := parser.parseUnaryExpression()

	return left
}

func (parser *Parser) parseUnaryExpression() ast.Expression {

	switch parser.token {
	case token.NOT:
		unaryExpression := &ast.UnaryExpression{
			Index:    parser.expect(parser.token),
			Operator: parser.token,
			Operand:  parser.parseUnaryExpression(),
		}
		return unaryExpression
	}

	left := parser.parseUpdateExpression()

	return left
}

func (parser *Parser) parseUpdateExpression() ast.Expression {
	isUpdate := true
	isUpdateToken := func(tkn token.Token) bool {
		return tkn == token.INCREMENT || tkn == token.DECREMENT
	}

	index := parser.index
	operator := parser.token
	var operand ast.Expression
	var isPostfix bool

	if isUpdateToken(operator) {
		isPostfix = false
	} else {
		operand = parser.parseLeftHandSideExpressionAllowCall()
		if isUpdateToken(parser.token) {
			isPostfix = true
			index = parser.index
			operator = parser.token
		} else {
			isUpdate = false
		}
	}

	if isUpdate {
		parser.next()
		switch operand.(type) {
		case *ast.Identifier:
		default:
			parser.error(index, "Invalid left-hand side in assignment")
			parser.nextStatement()
			return &ast.BadExpression{Start: index, End: parser.index}
		}
		return &ast.UnaryExpression{
			Index:    index,
			Operator: operator,
			Operand:  operand,
			Postfix:  isPostfix,
		}
	}
	return operand
}

func (parser *Parser) parseLeftHandSideExpressionAllowCall() ast.Expression {
	left := parser.parsePrimaryExpression()

	switch parser.token {
	case token.LEFT_PARENTHESIS:
		left = parser.parseCallExpression(left)
	}
	return left
}

func (parser *Parser) parsePrimaryExpression() ast.Expression {
	index := parser.index

	switch parser.token {
	case token.IDENTIFIER:
		return parser.parseIdentifier()
	case token.NUMBER:
		return parser.parseNumberLiteral()
	case token.STRING:
		return parser.parseStringLiteral()
	case token.BOOLEAN:
		return parser.parseBooleanLiteral()
	case token.NULL:
		return parser.parseNullLiteral()
	case token.LEFT_PARENTHESIS:
		return parser.parseParenthesisedExpression()
	}

	parser.errorUnexpectedToken(parser.token)
	parser.nextStatement()
	return &ast.BadExpression{
		Start: index,
		End:   parser.index,
	}
}

func (parser *Parser) parseNumberLiteral() ast.Expression {
	defer parser.expect(token.NUMBER)
	return &ast.NumberLiteral{
		Index:   parser.index,
		Literal: parser.literal,
		Value:   parser.parseNumberLiteralValue(parser.literal),
	}
}

func (parser *Parser) parseNumberLiteralValue(literal string) any {
	var value any = 0
	updateValue := func(v any, err error) bool {
		if err != nil {
			return false
		}
		value = v
		return true
	}
	intValue, err := strconv.ParseInt(parser.literal, 0, 64)
	if updateValue(intValue, err) {
		return value
	}
	floatValue, err := strconv.ParseFloat(parser.literal, 64)
	if updateValue(floatValue, err) {
		return value
	}
	return value
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	defer parser.expect(token.STRING)
	return &ast.StringLiteral{
		Index:   parser.index,
		Literal: parser.literal,
		Value:   parser.literal,
	}
}

func (parser *Parser) parseBooleanLiteral() ast.Expression {
	defer parser.expect(token.BOOLEAN)
	return &ast.BooleanLiteral{
		Index: parser.index,
		Value: parser.literal == "true",
	}
}
func (parser *Parser) parseNullLiteral() ast.Expression {
	defer parser.expect(token.NULL)
	return &ast.NullLiteral{
		Index: parser.index,
	}
}

func (parser *Parser) parseParenthesisedExpression() ast.Expression {
	parser.expect(token.LEFT_PARENTHESIS)
	left := parser.parseExpression()
	parser.expect(token.RIGHT_PARENTHESIS)
	return left
}

func (parser *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	leftParenthesis, arguments, rightParenthesis := parser.parseArguments()
	return &ast.CallExpression{
		Callee:           left,
		LeftParenthesis:  leftParenthesis,
		Arguments:        arguments,
		RightParenthesis: rightParenthesis,
	}
}

func (parser *Parser) parseArguments() (leftParenthesis file.Index, arguments []ast.Expression, rightParenthesis file.Index) {
	leftParenthesis = parser.expect(token.LEFT_PARENTHESIS)
	for parser.token != token.RIGHT_PARENTHESIS {
		arguments = append(arguments, parser.parseExpression())
		if parser.token != token.COMMA {
			break
		}
		parser.expect(token.COMMA)
	}
	rightParenthesis = parser.expect(token.RIGHT_PARENTHESIS)
	return
}
