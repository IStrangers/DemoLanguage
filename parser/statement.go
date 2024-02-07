package parser

import (
	"github.com/istrangers/demolanguage/ast"
	"github.com/istrangers/demolanguage/file"
	"github.com/istrangers/demolanguage/token"
)

func (parser *Parser) parseStatementList() []ast.Statement {
	return parser.parseStatementListByCondition(func(tkn token.Token) bool {
		return tkn != token.EOF
	})
}

func (parser *Parser) parseStatementListByCondition(endCondition func(token.Token) bool) []ast.Statement {
	var statementList []ast.Statement
	for endCondition(parser.token) {
		statementList = append(statementList, parser.parseStatement())
	}
	return statementList
}

func (parser *Parser) nextStatement() {
	for {
		switch parser.token {
		case token.VAR, token.FUN, token.RETURN, token.IF,
			token.FOR, token.SWITCH, token.BREAK, token.CONTINUE,
			token.THROW, token.TRY, token.CATCH, token.FINALLY:
			return
		case token.EOF:
			return
		}
		parser.next()
	}
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.token {
	case token.EOF:
		parser.errorUnexpectedToken(parser.token)
		return &ast.BadStatement{Start: parser.index, End: parser.index + 1}
	case token.LEFT_BRACE:
		return parser.parseBlockStatement()
	case token.VAR:
		return parser.parseVarStatement()
	case token.FUN:
		return parser.parseFunStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	case token.IF:
		return parser.parseIfStatement()
	case token.FOR:
		return parser.parseForStatement()
	case token.SWITCH:
		return parser.parseSwitchStatement()
	case token.BREAK:
		return parser.parseBreakStatement()
	case token.CONTINUE:
		return parser.parseContinueStatement()
	case token.THROW:
		return parser.parseThrowStatement()
	case token.TRY:
		return parser.parseTryCatchFinallyStatement()
	case token.CLASS:
		return parser.parseClassDeclaration()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	return &ast.BlockStatement{
		LeftBrace: parser.expect(token.LEFT_BRACE),
		Body: parser.parseStatementListByCondition(func(tkn token.Token) bool {
			return tkn != token.RIGHT_BRACE && tkn != token.EOF
		}),
		RightBrace: parser.expect(token.RIGHT_BRACE),
	}
}

func (parser *Parser) parseVarStatement() ast.Statement {
	varIndex := parser.expect(token.VAR)
	list := parser.parseVarDeclarationList(varIndex)
	return &ast.VarStatement{
		Var:  varIndex,
		List: list,
	}
}

func (parser *Parser) parseVarDeclarationList(varIndex file.Index) []*ast.Binding {
	bindingList := parser.parseBindingList()

	parser.scope.AddDeclaration(&ast.VariableDeclaration{
		Var:  varIndex,
		List: bindingList,
	})

	return bindingList
}

func (parser *Parser) parseFunStatement() ast.Statement {
	funStatement := &ast.FunStatement{
		FunLiteral: parser.parseFunLiteral(),
	}
	return funStatement
}

func (parser *Parser) parseReturnStatement() ast.Statement {
	returnIndex := parser.expect(token.RETURN)
	if !parser.scope.inFunction {
		parser.error(returnIndex, "Illegal return statement")
		parser.nextStatement()
		return &ast.BadStatement{Start: returnIndex, End: parser.index}
	}
	return &ast.ReturnStatement{
		Return:    returnIndex,
		Arguments: parser.parseReturnArguments(),
	}
}

func (parser *Parser) parseIfStatement() ast.Statement {
	parser.openScope()
	defer parser.closeScope()
	ifStatement := &ast.IfStatement{
		If:        parser.expect(token.IF),
		Condition: parser.parseExpression(),
	}
	ifStatement.Consequent = parser.parseBlockStatement()
	if parser.token == token.ELSE {
		parser.next()
		ifStatement.Alternate = parser.parseStatement()
	}
	return ifStatement
}

func (parser *Parser) parseForStatement() ast.Statement {
	parser.openScope()
	defer parser.closeScope()
	forStatement := &ast.ForStatement{
		For: parser.expect(token.FOR),
	}
	if parser.token != token.LEFT_BRACE {
		if parser.token != token.SEMICOLON {
			forStatement.Initializer = parser.parseVarStatement()
		}
		if parser.token == token.SEMICOLON {
			parser.expect(token.SEMICOLON)
		}
		if parser.token != token.SEMICOLON {
			forStatement.Condition = parser.parseExpression()
		}
		if parser.token == token.SEMICOLON {
			parser.expect(token.SEMICOLON)
		}
		if parser.token != token.SEMICOLON && parser.token != token.LEFT_BRACE {
			forStatement.Update = parser.parseExpression()
		}
	}
	parser.scope.inIteration = true
	forStatement.Body = parser.parseBlockStatement()
	parser.scope.inIteration = false
	return forStatement
}

func (parser *Parser) parseSwitchStatement() ast.Statement {
	switchStatement := &ast.SwitchStatement{
		Switch:       parser.expect(token.SWITCH),
		Discriminant: parser.parseExpression(),
		Default:      -1,
	}
	parser.scope.inSwitch = true
	switchStatement.Body, switchStatement.Default = parser.parseCaseStatementList()
	parser.scope.inSwitch = false
	switchStatement.RightBrace = parser.expect(token.RIGHT_BRACE)
	return switchStatement
}

func (parser *Parser) parseCaseStatementList() ([]*ast.CaseStatement, int) {
	var caseStatementList []*ast.CaseStatement
	var defaultIndex = -1
	parser.expect(token.LEFT_BRACE)
	for index := 0; parser.token != token.RIGHT_BRACE && parser.token != token.EOF; index++ {
		caseStatement := parser.parseCaseStatement()
		caseStatementList = append(caseStatementList, caseStatement)
		if caseStatement.Condition == nil {
			if defaultIndex == -1 {
				defaultIndex = index
			} else {
				parser.error(caseStatement.StartIndex(), "Already saw a default in switch")
			}
		}
	}
	return caseStatementList, defaultIndex
}

func (parser *Parser) parseCaseStatement() *ast.CaseStatement {
	parser.openScope()
	defer parser.closeScope()
	tkn := parser.token
	caseStatement := &ast.CaseStatement{
		Case:      parser.expect(parser.token),
		Condition: nil,
	}
	if tkn == token.CASE {
		caseStatement.Condition = parser.parseExpression()
	}
	caseStatement.Consequent = parser.parseStatement()
	return caseStatement
}

func (parser *Parser) parseBreakStatement() ast.Statement {
	breakIndex := parser.expect(token.BREAK)
	if !parser.scope.inIteration {
		parser.error(breakIndex, "Illegal break statement")
		parser.nextStatement()
		return &ast.BadStatement{Start: breakIndex, End: parser.index}
	}
	return &ast.BreakStatement{
		Break: breakIndex,
	}
}

func (parser *Parser) parseContinueStatement() ast.Statement {
	continueIndex := parser.expect(token.CONTINUE)
	if !parser.scope.inIteration {
		parser.error(continueIndex, "Illegal continue statement")
		parser.nextStatement()
		return &ast.BadStatement{Start: continueIndex, End: parser.index}
	}
	return &ast.ContinueStatement{
		Continue: continueIndex,
	}
}

func (parser *Parser) parseThrowStatement() ast.Statement {
	return &ast.ThrowStatement{
		Throw:    parser.expect(token.THROW),
		Argument: parser.parseReturnArguments()[0],
	}
}

func (parser *Parser) parseTryCatchFinallyStatement() ast.Statement {
	try := parser.expect(token.TRY)
	tryBody := parser.parseBlockStatement()
	statement := &ast.TryCatchFinallyStatement{
		Try:     try,
		TryBody: tryBody,
	}
	if parser.token == token.CATCH {
		parser.expect(token.CATCH)
		statement.CatchParameters = parser.parseParameterList()
		statement.CatchBody = parser.parseBlockStatement()
	}
	if parser.token == token.FINALLY {
		parser.expect(token.FINALLY)
		statement.FinallyBody = parser.parseBlockStatement()
	}
	return statement
}

func (parser *Parser) parseClassDeclaration() ast.Statement {
	classDeclaration := &ast.ClassDeclaration{
		Index: parser.expect(token.CLASS),
		Name:  parser.parseIdentifier(),
	}

	if parser.token == token.EXTENDS {
		parser.expect(token.EXTENDS)
		classDeclaration.SuperClass = parser.parseIdentifier()
	}

	if parser.token == token.IMPLEMENTS {
		parser.expect(token.IMPLEMENTS)
		for parser.token != token.LEFT_BRACE && parser.token != token.EOF {
			classDeclaration.Interfaces = append(classDeclaration.Interfaces, parser.parseIdentifier())
			if parser.token == token.COMMA {
				parser.expect(token.COMMA)
			} else {
				break
			}
		}
	}

	classDeclaration.LeftBrace = parser.expect(token.LEFT_BRACE)
	classDeclaration.Body = parser.parseDeclarations()
	classDeclaration.RightBrace = parser.expect(token.RIGHT_BRACE)
	classDeclaration.ClassDefinition = parser.slice(classDeclaration.StartIndex(), classDeclaration.EndIndex())
	return classDeclaration
}

func (parser *Parser) parseDeclarations() (declarations []ast.Declaration) {
	for parser.token != token.RIGHT_BRACE && parser.token != token.EOF {
		declarations = append(declarations, parser.parseDeclaration())
	}
	return
}

func (parser *Parser) parseDeclaration() ast.Declaration {
	index := parser.index
	switch parser.token {
	case token.STATIC:
		return &ast.StaticBlockDeclaration{
			Index: parser.expect(token.STATIC),
			Body:  parser.parseBlockStatement(),
		}
	case token.PRIVATE, token.PROTECTED, token.PUBLIC:
		index := parser.expect(parser.token)
		static := false
		if parser.token == token.STATIC {
			parser.expect(token.STATIC)
			static = true
		}
		name := parser.parseIdentifier()
		if parser.token == token.LEFT_PARENTHESIS {
			funLiteral := &ast.FunLiteral{
				Fun:  index,
				Name: name,
			}
			return &ast.MethodDeclaration{
				Index:          index,
				AccessModifier: parser.token,
				Static:         static,
				Body:           parser.parseAnonymousFunLiteral(funLiteral),
			}
		} else {
			fieldDeclaration := &ast.FieldDeclaration{
				Index:          index,
				AccessModifier: parser.token,
				Static:         static,
				Name:           name,
			}
			if parser.token == token.ASSIGN {
				parser.expect(token.ASSIGN)
				fieldDeclaration.Initializer = parser.parseExpression()
			}
			return fieldDeclaration
		}
	default:
		parser.error(parser.index, "Illegal break declaration")
		parser.nextStatement()
		return &ast.BadDeclaration{
			Start: index,
			End:   parser.index,
		}
	}
}

func (parser *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Expression: parser.parseExpression(),
	}
}
