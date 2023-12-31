package vm

import "github.com/istrangers/demolanguage/ast"

func (self *Compiler) checkStatementSyntax(st ast.Statement) {
	exitVirtualMode := self.enterVirtualMode()
	defer exitVirtualMode()
	self.compileStatement(st, false)
}

func (self *Compiler) isEmptyResultStatement(st ast.Statement) bool {
	switch st := st.(type) {
	case *ast.VarStatement, *ast.BreakStatement, *ast.ContinueStatement, *ast.FunStatement:
		return true
	case *ast.BlockStatement:
		for _, s := range st.Body {
			if _, ok := s.(*ast.BreakStatement); ok {
				return true
			}
			if _, ok := s.(*ast.ContinueStatement); ok {
				return true
			}
			if self.isEmptyResultStatement(s) {
				continue
			}
			return false
		}
		return true
	}
	return false
}

func (self *Compiler) compileStatements(statements []ast.Statement, needResult bool) {
	lastNeedResultIndex := -1
	for i, statement := range statements {
		if self.isEmptyResultStatement(statement) {
			continue
		}
		lastNeedResultIndex = i
	}
	for i, statement := range statements {
		if lastNeedResultIndex == i {
			self.compileStatement(statement, needResult)
		} else {
			self.compileStatement(statement, false)
		}
	}
}

func (self *Compiler) compileStatement(statement ast.Statement, needResult bool) {
	switch st := statement.(type) {
	case *ast.BlockStatement:
		self.compileBlockStatement(st, needResult)
	case *ast.VarStatement:
		self.compileVarStatement(st)
	case *ast.BreakStatement:
		self.compileBreakStatement(st)
	case *ast.ContinueStatement:
		self.compileContinueStatement(st)
	case *ast.ReturnStatement:
		self.compileReturnStatement(st)
	case *ast.IfStatement:
		self.compileIfStatement(st, needResult)
	case *ast.SwitchStatement:
		self.compileSwitchStatement(st, needResult)
	case *ast.ForStatement:
		self.compileForStatement(st, needResult)
	case *ast.FunStatement:
		self.compileFunStatement(st)
	case *ast.ExpressionStatement:
		self.compileExpressionStatement(st, needResult)
	}
}

func (self *Compiler) compileEnterBlockStatements(sts []ast.Statement, blockType BlockType, callback func()) {
	scopeDeclared := self.isScopeDeclared(sts)

	var enter *EnterBlock
	if scopeDeclared {
		self.openScopeNested()
		self.openBlock(blockType)
		self.compileScopeDeclared(sts)
		enter = &EnterBlock{}
		self.addProgramInstructions(enter)
	}

	callback()

	if scopeDeclared {
		self.leaveBlockScope(enter)
		self.closeScope()
	}
}

func (self *Compiler) compileBlockStatement(st *ast.BlockStatement, needResult bool) {
	self.compileEnterBlockStatements(st.Body, BlockScope, func() {
		self.compileStatements(st.Body, needResult)
	})
}

func (self *Compiler) compileVarStatement(st *ast.VarStatement) {
	for _, binding := range st.List {
		switch target := binding.Target.(type) {
		case *ast.Identifier:
			self.emitVarAssign(target.Name, int(target.StartIndex()-1), self.compileExpression(binding.Initializer))
		default:
			self.throwSyntaxError(int(target.StartIndex()-1), "unsupported variable binding target: %T", target)
		}
	}
}

func (self *Compiler) compileBreakStatement(st *ast.BreakStatement) {
	index := self.getInstructionSize()
	self.addProgramInstructions(nil)
	block := self.findBlockByType(BlockLoop)
	block.breaks = append(self.block.breaks, index)
}

func (self *Compiler) compileContinueStatement(st *ast.ContinueStatement) {
	index := self.getInstructionSize()
	self.addProgramInstructions(nil)
	block := self.findBlockByType(BlockLoop)
	block.continues = append(self.block.continues, index)
}

func (self *Compiler) compileReturnStatement(st *ast.ReturnStatement) {
	if st.Arguments != nil && len(st.Arguments) > 0 {
		for _, argument := range st.Arguments {
			self.chooseHandlingGetterExpression(self.compileExpression(argument), true)
		}
	} else {
		self.addProgramInstructions(LoadNull)
	}
	self.addProgramInstructions(Ret)
}

func (self *Compiler) compileIfStatement(st *ast.IfStatement, needResult bool) {
	conditionExpr := self.compileExpression(st.Condition)
	if conditionExpr.isConstExpression() {
		conditionValue := self.evalConstValueExpr(conditionExpr)
		if conditionValue.toBool() {
			self.compileStatement(st.Consequent, needResult)
			self.checkStatementSyntax(st.Alternate)
		} else if st.Alternate != nil {
			self.checkStatementSyntax(st.Consequent)
			self.compileStatement(st.Alternate, needResult)
		}
	} else {
		self.handlingGetterExpression(conditionExpr, true)
		consequentJmp := self.getInstructionSize()
		self.addProgramInstructions(nil)
		self.compileStatement(st.Consequent, needResult)
		if st.Alternate != nil {
			alternateJmp := self.getInstructionSize()
			self.addProgramInstructions(nil)
			self.setProgramInstruction(consequentJmp, Jne(self.getInstructionSize()-consequentJmp))
			self.compileStatement(st.Alternate, needResult)
			self.setProgramInstruction(alternateJmp, Jump(self.getInstructionSize()-alternateJmp))
		} else {
			self.setProgramInstruction(consequentJmp, Jne(self.getInstructionSize()-consequentJmp))
		}
	}
}

func (self *Compiler) compileSwitchStatement(st *ast.SwitchStatement, needResult bool) {
	discriminantExpr := self.compileExpression(st.Discriminant)
	if discriminantExpr.isConstExpression() {
		discriminantValue := self.evalConstValueExpr(discriminantExpr)
		var consequent ast.Statement
		for index, caseStatement := range st.Body {
			if index == st.Default && consequent == nil {
				consequent = caseStatement.Consequent
			} else {
				conditionExpr := self.compileExpression(caseStatement.Condition)
				if !conditionExpr.isConstExpression() {
					self.throwSyntaxError(int(caseStatement.StartIndex()-1), "Expression is not a constant")
					return
				}
				caseValue := self.evalConstValueExpr(conditionExpr)
				if discriminantValue.sameAs(caseValue) {
					consequent = caseStatement.Consequent
				}
			}
			self.checkStatementSyntax(caseStatement.Consequent)
		}
		if consequent == nil {
			return
		}
		self.compileStatement(consequent, needResult)
	} else {
		var jumpInstructionIndexs []int
		self.handlingGetterExpression(discriminantExpr, true)
		var defaultCaseStatement *ast.CaseStatement
		for index, caseStatement := range st.Body {
			if index == st.Default {
				defaultCaseStatement = caseStatement
				continue
			}
			conditionExpr := self.compileExpression(caseStatement.Condition)
			if conditionExpr.isConstExpression() {
				self.addProgramInstructions(Dup)
				self.chooseHandlingGetterExpression(conditionExpr, true)
				self.addProgramInstructions(EQ)
			} else {
				self.chooseHandlingGetterExpression(conditionExpr, true)
			}
			caseJmp := self.getInstructionSize()
			self.addProgramInstructions(nil)
			self.compileStatement(caseStatement.Consequent, needResult)
			jumpInstructionIndexs = append(jumpInstructionIndexs, self.getInstructionSize())
			self.addProgramInstructions(nil)
			self.setProgramInstruction(caseJmp, Jne(self.getInstructionSize()-caseJmp))
		}
		if defaultCaseStatement != nil {
			self.compileStatement(defaultCaseStatement.Consequent, needResult)
		}
		jump := self.getInstructionSize()
		for _, jumpInstructionIndex := range jumpInstructionIndexs {
			self.setProgramInstruction(jumpInstructionIndex, Jump(jump-jumpInstructionIndex))
		}
	}
}

func (self *Compiler) compileForStatement(st *ast.ForStatement, needResult bool) {
	blockLoop := self.openBlockLoop()

	self.compileEnterBlockStatements([]ast.Statement{st.Initializer}, BlockIterator, func() {
		if st.Initializer != nil {
			self.compileStatement(st.Initializer, true)
		}
		jumpIndex := self.getInstructionSize()
		conditionJumpIndex := -1
		if st.Condition != nil {
			conditionExpr := self.compileExpression(st.Condition)
			self.chooseHandlingGetterExpression(conditionExpr, true)
			conditionJumpIndex = self.getInstructionSize()
			self.addProgramInstructions(nil)
		}
		self.compileStatement(st.Body, needResult)
		blockLoop.continueBase = self.getInstructionSize()
		if st.Update != nil {
			updateExpr := self.compileExpression(st.Update)
			self.handlingGetterExpression(updateExpr, true)
		}
		self.addProgramInstructions(Jump(jumpIndex - self.getInstructionSize()))
		if conditionJumpIndex != -1 {
			self.setProgramInstruction(conditionJumpIndex, Jne(self.getInstructionSize()-conditionJumpIndex))
		}
	})

	self.closeBlock()
}

func (self *Compiler) compileFunStatement(st *ast.FunStatement) {
	funLiteralExpr := self.compileExpression(st.FunLiteral)
	self.handlingGetterExpression(funLiteralExpr, true)
}

func (self *Compiler) compileExpressionStatement(st *ast.ExpressionStatement, needResult bool) {
	self.chooseHandlingGetterExpression(self.compileExpression(st.Expression), needResult)
	if !needResult {
		return
	}
	self.addProgramInstructions(SaveResult)
}
