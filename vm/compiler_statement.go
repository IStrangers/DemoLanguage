package vm

import "DemoLanguage/ast"

func (self *Compiler) checkStatementSyntax(st ast.Statement) {
	exitVirtualMode := self.enterVirtualMode()
	defer exitVirtualMode()
	self.compileStatement(st, false)
}

func (self *Compiler) compileStatements(statements []ast.Statement, needResult bool) {
	for _, statement := range statements {
		self.compileStatement(statement, needResult)
	}
}

func (self *Compiler) compileStatement(statement ast.Statement, needResult bool) {
	switch st := statement.(type) {
	case *ast.BlockStatement:
		self.compileBlockStatement(st, needResult)
	case *ast.VarStatement:
		self.compileVarStatement(st)
	case *ast.IfStatement:
		self.compileIfStatement(st, needResult)
	case *ast.SwitchStatement:
		self.compileSwitchStatement(st, needResult)
	case *ast.ForStatement:
		self.compileForStatement(st, needResult)
	case *ast.ExpressionStatement:
		self.compileExpressionStatement(st, needResult)
	}
}

func (self *Compiler) compileBlockStatement(st *ast.BlockStatement, needResult bool) {
	self.compileStatements(st.Body, needResult)
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
		program := self.program
		consequentJmp := program.getInstructionSize()
		self.addProgramInstructions(nil)
		self.compileStatement(st.Consequent, needResult)
		if st.Alternate != nil {
			alternateJmp := program.getInstructionSize()
			self.addProgramInstructions(nil)
			self.setProgramInstruction(consequentJmp, Jne(program.getInstructionSize()-consequentJmp))
			self.compileStatement(st.Alternate, needResult)
			self.setProgramInstruction(alternateJmp, Jump(program.getInstructionSize()-alternateJmp))
		} else {
			self.setProgramInstruction(consequentJmp, Jne(program.getInstructionSize()-consequentJmp))
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
		program := self.program
		self.handlingGetterExpression(discriminantExpr, needResult)
		var defaultCaseStatement *ast.CaseStatement
		for index, caseStatement := range st.Body {
			if index == st.Default {
				defaultCaseStatement = caseStatement
				continue
			}
			conditionExpr := self.compileExpression(caseStatement.Condition)
			if conditionExpr.isConstExpression() {
				self.addProgramInstructions(Dup)
				self.chooseHandlingGetterExpression(conditionExpr, needResult)
				self.addProgramInstructions(EQ)
			} else {
				self.chooseHandlingGetterExpression(conditionExpr, needResult)
			}
			caseJmp := self.program.getInstructionSize()
			self.addProgramInstructions(nil)
			self.compileStatement(caseStatement.Consequent, needResult)
			jumpInstructionIndexs = append(jumpInstructionIndexs, program.getInstructionSize())
			self.addProgramInstructions(nil)
			self.setProgramInstruction(caseJmp, Jne(program.getInstructionSize()-caseJmp))
		}
		if defaultCaseStatement != nil {
			self.compileStatement(defaultCaseStatement.Consequent, needResult)
		}
		jump := program.getInstructionSize()
		for _, jumpInstructionIndex := range jumpInstructionIndexs {
			self.setProgramInstruction(jumpInstructionIndex, Jump(jump-jumpInstructionIndex))
		}
	}
}

func (self *Compiler) compileForStatement(st *ast.ForStatement, needResult bool) {
	program := self.program
	if st.Initializer != nil {
		self.compileStatement(st.Initializer, needResult)
	}
	jumpIndex := program.getInstructionSize()
	conditionJumpIndex := -1
	if st.Condition != nil {
		conditionExpr := self.compileExpression(st.Condition)
		self.chooseHandlingGetterExpression(conditionExpr, needResult)
		conditionJumpIndex = program.getInstructionSize()
		self.addProgramInstructions(nil)
	}
	self.compileStatement(st.Body, needResult)
	if st.Update != nil {
		updateExpr := self.compileExpression(st.Update)
		self.handlingGetterExpression(updateExpr, needResult)
	}
	self.addProgramInstructions(Jump(jumpIndex - program.getInstructionSize()))
	if conditionJumpIndex != -1 {
		self.setProgramInstruction(conditionJumpIndex, Jne(program.getInstructionSize()-conditionJumpIndex))
	}
}

func (self *Compiler) compileExpressionStatement(st *ast.ExpressionStatement, needResult bool) {
	self.chooseHandlingGetterExpression(self.compileExpression(st.Expression), needResult)
	if !needResult {
		return
	}
	self.addProgramInstructions(SaveResult)
}
