package vm

import (
	"github.com/istrangers/demolanguage/ast"
	"github.com/istrangers/demolanguage/token"
)

func (self *Compiler) checkStatementSyntax(st ast.Statement) {
	exitVirtualMode := self.enterVirtualMode()
	defer exitVirtualMode()
	self.compileStatement(st, false)
}

func (self *Compiler) isEmptyResultStatement(st ast.Statement) bool {
	switch st := st.(type) {
	case *ast.VarStatement, *ast.BreakStatement, *ast.ContinueStatement, *ast.FunStatement, *ast.ClassDeclaration:
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

func (self *Compiler) findLastResultIndex(statements []ast.Statement) (lastResultIndex int, breakingBlock *Block) {
	lastResultIndex = -1
	for i, statement := range statements {
		if st, ok := statement.(ast.BranchStatement); ok {
			if block := self.findBlockByType([]BlockType{BlockLoop, BlockSwitch}, st.Token() == token.BREAK); block != nil {
				breakingBlock = block
			}
			break
		}
		if self.isEmptyResultStatement(statement) {
			continue
		}
		lastResultIndex = i
	}
	return
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
	case ast.BranchStatement:
		self.compileBranchStatement(st)
	case *ast.ReturnStatement:
		self.compileReturnStatement(st)
	case *ast.IfStatement:
		self.compileIfStatement(st, needResult)
	case *ast.SwitchStatement:
		self.compileSwitchStatement(st, needResult)
	case *ast.ForStatement:
		self.compileForStatement(st, needResult)
	case *ast.ThrowStatement:
		self.compileThrowStatement(st)
	case *ast.TryCatchFinallyStatement:
		self.compileTryCatchFinallyStatement(st, needResult)
	case *ast.FunStatement:
		self.compileFunStatement(st)
	case *ast.ClassDeclaration:
		self.compileClassDeclaration(st, needResult)
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
	block := self.findBlockByType([]BlockType{BlockLoop, BlockSwitch}, true)
	block.breaks = append(self.block.breaks, index)
}

func (self *Compiler) compileContinueStatement(st *ast.ContinueStatement) {
	index := self.getInstructionSize()
	self.addProgramInstructions(nil)
	block := self.findBlockByType([]BlockType{BlockLoop}, false)
	block.continues = append(self.block.continues, index)
}

func (self *Compiler) compileBranchStatement(st ast.BranchStatement) {
	if st.Token() == token.BREAK {
		self.compileBreakStatement(st.(*ast.BreakStatement))
	} else {
		self.compileContinueStatement(st.(*ast.ContinueStatement))
	}
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

func (self *Compiler) compileThrowStatement(st *ast.ThrowStatement) {
	expr := self.compileExpression(st.Argument)
	self.handlingGetterExpression(expr, true)
	expr.addSourceMap()
	self.addProgramInstructions(Throw)
}

func (self *Compiler) compileTryCatchFinallyStatement(st *ast.TryCatchFinallyStatement, needResult bool) {
	self.block = self.openBlockTry()
	var bodyNeedResult bool = needResult
	//var lp int
	//var finallyBreaking *Block
	//if st.FinallyBody != nil {
	//	lp, finallyBreaking = self.findLastResultIndex(st.FinallyBody.Body)
	//}
	//if finallyBreaking != nil {
	//	self.block.breaking = finallyBreaking
	//	if lp == -1 {
	//		bodyNeedResult = finallyBreaking.needResult
	//	}
	//}
	enterTryIndex := self.getInstructionSize()
	self.addProgramInstructions(nil)
	self.compileBlockStatement(st.TryBody, bodyNeedResult)

	var catchOffset int
	if st.CatchBody != nil {
		jumpIndex := self.getInstructionSize()
		self.addProgramInstructions(nil)
		catchOffset = self.getInstructionSize() - enterTryIndex
		if st.CatchParameters != nil {
			self.block = self.openBlockScope()
			self.openScopeNested()
			for _, binding := range st.CatchParameters.List {
				self.scope.bindName(binding.Target.(*ast.Identifier).Name)
			}
			enterBlock := &EnterBlock{}
			self.addProgramInstructions(enterBlock)
			self.compileStatements(st.CatchBody.Body, bodyNeedResult)
			self.leaveBlockScope(enterBlock)
			if self.scope.bindings[0].inStash {
				self.setProgramInstruction(enterTryIndex+catchOffset, &EnterCatchBlock{
					stackSize: enterBlock.stackSize,
					stashSize: enterBlock.stashSize,
				})
			} else {
				enterBlock.stackSize--
			}
			self.closeScope()
		} else {
			self.addProgramInstructions(Pop)
			self.compileBlockStatement(st.CatchBody, bodyNeedResult)
		}
		self.setProgramInstruction(jumpIndex, Jump(self.getInstructionSize()-jumpIndex))
	}

	var finallyOffset int
	if st.FinallyBody != nil {
		self.addProgramInstructions(EnterFinally{})
		finallyOffset = self.getInstructionSize() - enterTryIndex
		self.compileBlockStatement(st.FinallyBody, false)
		self.addProgramInstructions(LeaveFinally{})
	} else {
		self.addProgramInstructions(LeaveTry{})
	}
	self.setProgramInstruction(enterTryIndex, EnterTry{catchOffset, finallyOffset})
	self.closeBlock()
}

func (self *Compiler) compileFunStatement(st *ast.FunStatement) {
	funLiteralExpr := self.compileExpression(st.FunLiteral)
	self.handlingGetterExpression(funLiteralExpr, true)
}

func (self *Compiler) compileClassDeclaration(st *ast.ClassDeclaration, needResult bool) {
	self.openScopeNested()

	enterBlock := &EnterBlock{}
	//markIndex := self.getInstructionSize()
	self.addProgramInstructions(enterBlock)
	self.block = self.openBlockScope()
	self.scope.bindName(st.Name.Name)
	self.openClassScope()

	newClass := &NewClass{
		name:   st.Name.Name,
		source: st.ClassDefinition,
	}

	var newClassInstruction Instruction
	isDerivedClass := false
	if st.SuperClass != nil {
		isDerivedClass = true
		newClassInstruction = &NewDerivedClass{
			newClass: newClass,
		}
	} else {
		newClassInstruction = newClass
	}

	var staticBlocks []*ast.StaticBlockDeclaration
	var instanceFieldDecls, staticFieldDecls []*ast.FieldDeclaration
	var instanceMethodDecls, staticMethodDecls []*ast.MethodDeclaration
	instanceCount, staticCount := 0, 0
	for _, declaration := range st.Body {
		switch decl := declaration.(type) {
		case *ast.StaticBlockDeclaration:
			if len(decl.Body.Body) > 0 {
				staticBlocks = append(staticBlocks, decl)
				staticCount++
			}
		case *ast.FieldDeclaration:
			if decl.Static {
				staticFieldDecls = append(staticFieldDecls, decl)
				staticCount++
			} else {
				instanceFieldDecls = append(instanceFieldDecls, decl)
				instanceCount++
			}
		case *ast.MethodDeclaration:
			if newClass.name == decl.Body.Name.Name {
				program, paramNum := self.compileConstructor(decl.Body, isDerivedClass)
				newClass.constructorList = append(newClass.constructorList, &Constructor{
					paramNum,
					program,
				})
			} else {
				if decl.Static {
					staticMethodDecls = append(staticMethodDecls, decl)
					staticCount++
				} else {
					instanceMethodDecls = append(instanceMethodDecls, decl)
					instanceCount++
				}
			}
		}
	}

	if instanceCount > 0 {
		newClass.init = self.compileDeclarations("<instance_members_initializer>", nil, instanceFieldDecls, instanceMethodDecls)
	}

	if isDerivedClass {
		self.handlingGetterExpression(self.compileExpression(st.SuperClass), true)
	}
	self.program.addSourceMap(int(st.StartIndex() - 1))
	self.addProgramInstructions(newClassInstruction)

	if staticCount > 0 {
		self.addProgramInstructions(&ClassStaticPropInit{
			init: self.compileDeclarations("<static_initializer>", staticBlocks, staticFieldDecls, staticMethodDecls),
		})
	}

	self.closeClassScope()
	self.closeScope()
}

func (self *Compiler) compileDeclarations(functionName string, blocks []*ast.StaticBlockDeclaration, fields []*ast.FieldDeclaration, methods []*ast.MethodDeclaration) *Program {
	originBlock, originProgram := self.block, self.program
	defer func() {
		self.block = originBlock
		self.program = originProgram
	}()

	self.openBlockScope()
	self.program = &Program{
		source:       originProgram.source,
		functionName: functionName,
		instructions: InstructionArray{},
	}

	self.openScope()
	for _, block := range blocks {
		self.addProgramInstructions(Dup)
		self.handlingGetterCompiledFunLiteralExpression(self.compileFunLiteral(&ast.FunLiteral{
			Fun:             block.Index,
			ParameterList:   &ast.ParameterList{},
			Body:            block.Body,
			DeclarationList: []*ast.VariableDeclaration{},
			FunDefinition:   block.Source,
		}), true)
		self.program.addSourceMap(int(block.Index - 1))
		self.addProgramInstructions(Call(0), Pop)
	}

	for _, field := range fields {
		if field.Initializer != nil {
			valueExpr := self.compileExpression(field.Initializer)
			self.chooseHandlingGetterExpression(valueExpr, true)
		} else {
			self.addProgramInstructions(LoadNull)
		}
		self.addProgramInstructions(AddProp(field.Name.Name))
	}

	for _, method := range methods {
		funLiteral := method.Body
		self.handlingGetterCompiledFunLiteralExpression(self.compileFunLiteral(funLiteral), true)
		self.addProgramInstructions(AddProp(funLiteral.Name.Name))
	}
	stashSize, stackSize := self.scope.finaliseVarAlloc(0)
	if stackSize != 0 {
		panic("Compiler bug: stackSize != 0 in initFields")
	}
	if stashSize > 0 {
		if stashSize != 1 {
			panic("Compiler bug: stashSize != 1 in initFields")
		}
		self.setProgramInstruction(0, &EnterFunStash{
			stashSize: 1,
		})
	}

	program := self.program
	self.closeScope()
	return program
}

func (self *Compiler) compileConstructor(funLiteral *ast.FunLiteral, isDerivedClass bool) (*Program, int) {
	return self.compileFunProgram(self.compileFunLiteral(funLiteral))
}

func (self *Compiler) compileExpressionStatement(st *ast.ExpressionStatement, needResult bool) {
	self.chooseHandlingGetterExpression(self.compileExpression(st.Expression), needResult)
	if !needResult {
		return
	}
	self.addProgramInstructions(SaveResult)
}
