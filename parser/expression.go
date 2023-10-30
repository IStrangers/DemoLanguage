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
	funLiteral.Fun = parser.expect(token.FUN)
	if parser.token != token.LEFT_PARENTHESIS {
		funLiteral.Name = parser.parseIdentifier()
	}
	funLiteral.ParameterList = parser.parseFunParameterList()
	funLiteral.Body, funLiteral.DeclarationList = parser.parseFunBlock(true)
	funLiteral.FunDefinition = parser.slice(funLiteral.StartIndex(), funLiteral.EndIndex())
	return funLiteral
}

func (parser *Parser) parseFunParameterList() *ast.ParameterList {
	parameterList := &ast.ParameterList{
		LeftParenthesis: parser.expect(token.LEFT_PARENTHESIS),
	}
	if parser.token != token.RIGHT_PARENTHESIS {
		parameterList.List = parser.parseBindingList()
	}
	parameterList.RightParenthesis = parser.expect(token.RIGHT_PARENTHESIS)
	return parameterList
}

func (parser *Parser) parseFunBlock(openScope bool) (ast.Statement, []*ast.VariableDeclaration) {
	if openScope {
		parser.openScope()
		defer parser.closeScope()
	}
	return parser.parseBlockStatement(), parser.scope.declarationList
}

func (parser *Parser) parseReturnArguments() (arguments []ast.Expression) {
	for parser.token != token.RIGHT_BRACE && parser.token != token.EOF {
		arguments = append(arguments, parser.parseExpression())
		if parser.token == token.COMMA {
			parser.expect(token.COMMA)
		} else {
			break
		}
	}
	return
}

func (parser *Parser) parseExpression() ast.Expression {
	left := parser.parseAssignExpression()

	return left
}

func (parser *Parser) parseAssignExpression() ast.Expression {
	parseState := parser.markParseState()
	parenthesis := false

	switch parser.token {
	case token.LEFT_PARENTHESIS:
		parenthesis = true
	}

	left := parser.parseConditionalExpression()

	var operator token.Token
	switch parser.token {
	case token.ASSIGN:
		operator = parser.token
	case token.ADDITION_ASSIGN:
		operator = token.ADDITION
	case token.SUBTRACT_ASSIGN:
		operator = token.SUBTRACT
	case token.MULTIPLY_ASSIGN:
		operator = token.MULTIPLY
	case token.DIVIDE_ASSIGN:
		operator = token.DIVIDE
	case token.REMAINDER_ASSIGN:
		operator = token.REMAINDER
	case token.AND_ARITHMETIC_ASSIGN:
		operator = token.AND_ARITHMETIC
	case token.OR_ARITHMETIC_ASSIGN:
		operator = token.OR_ARITHMETIC
	case token.ARROW:
		parser.restoreParseState(parseState)
		left = parser.parseArrowFunctionLiteral()
	}

	if operator != 0 {
		index := parser.index
		err := true

		switch left.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
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
		parser.expect(parser.token)
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

	for {
		switch parser.token {
		case token.LOGICAL_OR:
			left = &ast.BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseLogicalAndExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseLogicalAndExpression() ast.Expression {
	left := parser.parseBitwiseOrExpression()

	for {
		switch parser.token {
		case token.LOGICAL_AND:
			left = &ast.BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseBitwiseOrExpression(),
			}
		default:
			return left
		}
	}
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

	for {
		switch parser.token {
		case token.EQUAL, token.NOT_EQUAL:
			left = &ast.BinaryExpression{
				Operator:   parser.expectToken(parser.token),
				Left:       left,
				Right:      parser.parseRelationalExpression(),
				Comparison: true,
			}
		default:
			return left
		}
	}
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

	for {
		switch parser.token {
		case token.ADDITION, token.SUBTRACT:
			left = &ast.BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseMultiplicativeExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseMultiplicativeExpression() ast.Expression {
	left := parser.parseExponentiationExpression()

	for {
		switch parser.token {
		case token.MULTIPLY, token.DIVIDE, token.REMAINDER:
			left = &ast.BinaryExpression{
				Operator: parser.expectToken(parser.token),
				Left:     left,
				Right:    parser.parseExponentiationExpression(),
			}
		default:
			return left
		}
	}
}

func (parser *Parser) parseExponentiationExpression() ast.Expression {
	left := parser.parseUnaryExpression()

	return left
}

func (parser *Parser) parseUnaryExpression() ast.Expression {

	tkn := parser.token
	switch tkn {
	case token.NOT, token.ADDITION, token.SUBTRACT:
		unaryExpression := &ast.UnaryExpression{
			Index:    parser.expect(tkn),
			Operator: tkn,
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
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
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

	for {
		switch parser.token {
		case token.DOT:
			left = parser.parseDotExpression(left)
			continue
		case token.LEFT_BRACKET:
			left = parser.parseBracketExpression(left)
			continue
		case token.LEFT_PARENTHESIS:
			left = parser.parseCallExpression(left)
			continue
		}
		break
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
	case token.LEFT_BRACKET:
		return parser.parseArrayLiteral()
	case token.LEFT_BRACE:
		return parser.parseObjectLiteral()
	case token.LEFT_PARENTHESIS:
		return parser.parseParenthesisedExpression()
	case token.THIS:
		return parser.parseThisExpression()
	case token.FUN:
		return parser.parseFunLiteral()
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
	intValue, err := strconv.ParseInt(literal, 0, 64)
	if updateValue(intValue, err) {
		return value
	}
	floatValue, err := strconv.ParseFloat(literal, 64)
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

func (parser *Parser) parseArrayLiteral() ast.Expression {
	arrayLiteral := &ast.ArrayLiteral{
		LeftBracket: parser.expect(token.LEFT_BRACKET),
	}
	var values []ast.Expression
	for parser.token != token.RIGHT_BRACKET && parser.token != token.EOF {
		if parser.token == token.COMMA {
			values = append(values, &ast.NullLiteral{Index: parser.index})
		} else {
			values = append(values, parser.parseExpression())
		}
		if parser.token != token.RIGHT_BRACKET {
			parser.expect(token.COMMA)
		}
	}
	arrayLiteral.Values = values
	arrayLiteral.RightBracket = parser.expect(token.RIGHT_BRACKET)
	return arrayLiteral
}

func (parser *Parser) parseObjectLiteral() ast.Expression {
	objectLiteral := &ast.ObjectLiteral{
		LeftBrace: parser.expect(token.LEFT_BRACE),
	}
	var properties []ast.Property
	for parser.token != token.RIGHT_BRACE && parser.token != token.EOF {
		property := parser.parseObjectProperty()
		if property != nil {
			properties = append(properties, property)
		}
		if parser.token != token.RIGHT_BRACE {
			parser.expect(token.COMMA)
		}
	}
	objectLiteral.Properties = properties
	objectLiteral.RightBrace = parser.expect(token.RIGHT_BRACE)
	return objectLiteral
}

func (parser *Parser) parseObjectProperty() ast.Property {
	propertyKeyValue := &ast.PropertyKeyValue{
		Name:  parser.parseIdentifier(),
		Colon: parser.expect(token.COLON),
		Value: parser.parseExpression(),
	}
	return propertyKeyValue
}

func (parser *Parser) parseArrowFunctionLiteral() ast.Expression {
	arrowFunctionLiteral := &ast.ArrowFunctionLiteral{
		Index: parser.index,
	}
	if parser.token == token.LEFT_PARENTHESIS {
		arrowFunctionLiteral.ParameterList = parser.parseFunParameterList()
	} else {
		identifier := parser.parseIdentifier()
		arrowFunctionLiteral.ParameterList = &ast.ParameterList{
			LeftParenthesis: identifier.StartIndex(),
			List: []*ast.Binding{{
				Target: identifier,
			}},
			RightParenthesis: identifier.EndIndex() - 1,
		}
	}
	arrowFunctionLiteral.Arrow = parser.expect(token.ARROW)
	if parser.token == token.LEFT_BRACE {
		arrowFunctionLiteral.Body, arrowFunctionLiteral.DeclarationList = parser.parseFunBlock(false)
	} else {
		expression := parser.parseExpression()
		arrowFunctionLiteral.Body = &ast.BlockStatement{
			LeftBrace:  expression.StartIndex(),
			Body:       []ast.Statement{&ast.ExpressionStatement{Expression: expression}},
			RightBrace: expression.EndIndex() - 1,
		}
	}
	arrowFunctionLiteral.FunDefinition = parser.slice(arrowFunctionLiteral.StartIndex(), arrowFunctionLiteral.EndIndex())
	return arrowFunctionLiteral
}

func (parser *Parser) parseParenthesisedExpression() ast.Expression {
	parser.expect(token.LEFT_PARENTHESIS)
	left := parser.parseExpression()
	parser.expect(token.RIGHT_PARENTHESIS)
	return left
}

func (parser *Parser) parseThisExpression() ast.Expression {
	defer parser.expect(token.THIS)
	return &ast.ThisExpression{
		Index: parser.index,
	}
}

func (parser *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	dotExpression := &ast.DotExpression{
		Left:       left,
		Dot:        parser.expect(token.DOT),
		Identifier: parser.parseIdentifier(),
	}
	return dotExpression
}

func (parser *Parser) parseBracketExpression(left ast.Expression) ast.Expression {
	dotExpression := &ast.BracketExpression{
		Left:         left,
		LeftBracket:  parser.expect(token.LEFT_BRACKET),
		Expression:   parser.parseExpression(),
		RightBracket: parser.expect(token.RIGHT_BRACKET),
	}
	return dotExpression
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
