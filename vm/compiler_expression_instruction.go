package vm

import (
	"DemoLanguage/ast"
	"DemoLanguage/token"
)

func (self *Compiler) evalConstExpr(expr CompiledExpression) (Value, *Exception) {
	if expr.isLiteralExpression() {
		return expr.(*CompiledLiteralExpression).value, nil
	}
	evalVM := self.evalVM
	originProgram := self.program
	isNewProgram := false
	if evalVM.program == nil {
		evalVM.program = &Program{
			source: self.program.source,
		}
		self.program = evalVM.program
		isNewProgram = true
	}
	savedPC := self.getInstructionSize()
	self.handlingGetterExpression(expr, true)
	evalVM.pc = savedPC
	ex := evalVM.runTry()
	if isNewProgram {
		evalVM.program = nil
		evalVM.pc = 0
		self.program = originProgram
	} else {
		evalVM.program.instructions = evalVM.program.instructions[:savedPC]
		self.program.instructions = evalVM.program.instructions
	}
	if ex != nil {
		return nil, ex
	}
	return evalVM.pop(), nil
}

func (self *Compiler) evalConstValueExpr(expr CompiledExpression) Value {
	value, ex := self.evalConstExpr(expr)
	if ex != nil {
		expr.addSourceMap()
		self.emitThrow(ex.value)
	}
	return value
}

func (self *Compiler) emitLoadValue(value Value, putOnStack bool) {
	if !putOnStack {
		return
	}
	index := self.addProgramValue(value)
	self.addProgramInstructions(LoadVal(index))
}

func (self *Compiler) emitVarAssign(name string, pos int, expr CompiledExpression) {
	if expr == nil {
		return
	}
	binding, exists := self.scope.lookupName(name)
	if exists {
		self.chooseHandlingGetterExpression(expr, true)
		self.program.addSourceMap(pos)
		binding.markAccessPoint(self.scope)
		self.addProgramInstructions(InitStackVar(0))
	} else {
		self.addProgramInstructions(ResolveVar(name))
		self.chooseHandlingGetterExpression(expr, true)
		self.program.addSourceMap(pos)
		self.addProgramInstructions(InitVar)
	}
}

func (self *Compiler) emitThrow(value Value) {

}

func (self *Compiler) handlingConstExpression(expr CompiledExpression, putOnStack bool) {
	value, ex := self.evalConstExpr(expr)
	if ex != nil {
		self.emitThrow(ex.value)
	} else {
		self.emitLoadValue(value, putOnStack)
	}
}

func (self *Compiler) handlingGetterExpression(expr CompiledExpression, putOnStack bool) {
	switch expr := expr.(type) {
	case *CompiledLiteralExpression:
		self.handlingGetterCompiledLiteralExpression(expr, putOnStack)
	case *CompiledObjectLiteralExpression:
		self.handlingGetterCompiledObjectLiteralExpression(expr, putOnStack)
	case *CompiledArrayLiteralExpression:
		self.handlingGetterCompiledArrayLiteralExpression(expr, putOnStack)
	case *CompiledIdentifierExpression:
		self.handlingGetterCompiledIdentifierExpression(expr, putOnStack)
	case *CompiledThisExpression:
		self.handlingGetterCompiledThisExpression(expr, putOnStack)
	case *CompiledUnaryExpression:
		self.handlingGetterCompiledUnaryExpression(expr, putOnStack)
	case *CompiledBinaryExpression:
		self.handlingGetterCompiledBinaryExpression(expr, putOnStack)
	case *CompiledAssignExpression:
		self.handlingGetterCompiledAssignExpression(expr, putOnStack)
	case *CompiledFunLiteralExpression:
		self.handlingGetterCompiledFunLiteralExpression(expr, putOnStack)
	case *CompiledCallExpression:
		self.handlingGetterCompiledCallExpression(expr, putOnStack)
	case *CompiledDotExpression:
		self.handlingGetterCompiledDotExpression(expr, putOnStack)
	case *CompiledBracketExpression:
		self.handlingGetterCompiledBracketExpression(expr, putOnStack)
	}
}

func (self *Compiler) chooseHandlingGetterExpression(expr CompiledExpression, putOnStack bool) {
	if expr.isConstExpression() {
		self.handlingConstExpression(expr, putOnStack)
	} else {
		self.handlingGetterExpression(expr, putOnStack)
	}
}

func (self *Compiler) handlingGetterCompiledLiteralExpression(expr *CompiledLiteralExpression, putOnStack bool) {
	self.emitLoadValue(expr.value, putOnStack)
}

func (self *Compiler) handlingGetterCompiledObjectLiteralExpression(expr *CompiledObjectLiteralExpression, putOnStack bool) {
	expr.addSourceMap()

	self.addProgramInstructions(NewObject)
	for _, property := range expr.properties {
		switch prop := property.(type) {
		case *ast.PropertyKeyValue:
			name := prop.Name.Name
			valueExpr := self.compileExpression(prop.Value)
			self.chooseHandlingGetterExpression(valueExpr, true)
			self.addProgramInstructions(AddProp(name))
		default:
			//wait adjust
			self.throwSyntaxError(expr.offset, "unknown Property type: %T", prop)
		}
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledArrayLiteralExpression(expr *CompiledArrayLiteralExpression, putOnStack bool) {
	expr.addSourceMap()
	self.addProgramInstructions(NewArray(len(expr.values)))
	for _, value := range expr.values {
		if value == nil {
			self.addProgramInstructions(LoadNull)
		} else {
			self.chooseHandlingGetterExpression(self.compileExpression(value), true)
		}
		self.addProgramInstructions(PushArrayValue)
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledIdentifierExpression(expr *CompiledIdentifierExpression, putOnStack bool) {
	expr.addSourceMap()

	binding, exists := self.scope.lookupName(expr.name)
	if exists {
		binding.markAccessPoint(self.scope)
		self.addProgramInstructions(LoadStackVar(0))
	} else {
		self.addProgramInstructions(LoadVar(expr.name))
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledThisExpression(expr *CompiledThisExpression, putOnStack bool) {
	expr.addSourceMap()

	binding, exists := self.scope.lookupName(thisBindingName)
	if exists {
		binding.markAccessPoint(self.scope)
		self.addProgramInstructions(LoadStackVar(0))
	} else {
		self.addProgramInstructions(LoadDynamicThis)
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledUnaryExpression(expr *CompiledUnaryExpression, putOnStack bool) {
	switch expr.operator {
	case token.NOT:
		self.chooseHandlingGetterExpression(expr.operand, true)
		self.addProgramInstructions(Not)
	case token.SUBTRACT:
		self.chooseHandlingGetterExpression(expr.operand, true)
		self.addProgramInstructions(Neg)
	case token.ADDITION:
		self.chooseHandlingGetterExpression(expr.operand, true)
	case token.INCREMENT:
		self.handlingUnaryExpression(expr.operand, func() {
			self.addProgramInstructions(Inc)
		}, expr.postfix, putOnStack)
	case token.DECREMENT:
		self.handlingUnaryExpression(expr.operand, func() {
			self.addProgramInstructions(Dec)
		}, expr.postfix, putOnStack)
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledBinaryExpression(expr *CompiledBinaryExpression, putOnStack bool) {
	operator := expr.operator
	if operator == token.LOGICAL_OR || operator == token.LOGICAL_AND {

		if expr.left.isConstExpression() {
			if v, ex := self.evalConstExpr(expr.left); ex == nil {
				boolVal := v.toBool()
				if (operator == token.LOGICAL_OR && boolVal) || (operator == token.LOGICAL_AND && !boolVal) {
					index := self.addProgramValue(v)
					self.addProgramInstructions(LoadVal(index))
				} else {
					self.chooseHandlingGetterExpression(expr.right, putOnStack)
				}
			} else {
				self.emitThrow(ex.value)
			}
		} else {
			self.handlingGetterExpression(expr.left, true)
			expr.addSourceMap()
			index := self.getInstructionSize()
			self.addProgramInstructions(nil)
			self.handlingGetterExpression(expr.right, true)
			var instruction Instruction
			switch operator {
			case token.LOGICAL_OR:
				instruction = Jeq1(self.getInstructionSize() - index)
			case token.LOGICAL_AND:
				instruction = Jne1(self.getInstructionSize() - index)
			}
			self.setProgramInstruction(index, instruction)
		}

	} else {

		self.chooseHandlingGetterExpression(expr.left, true)
		self.chooseHandlingGetterExpression(expr.right, true)
		expr.addSourceMap()

		switch expr.operator {
		case token.ADDITION:
			self.addProgramInstructions(Add)
		case token.SUBTRACT:
			self.addProgramInstructions(Sub)
		case token.MULTIPLY:
			self.addProgramInstructions(Mul)
		case token.DIVIDE:
			self.addProgramInstructions(Div)
		case token.REMAINDER:
			self.addProgramInstructions(Mod)
		case token.EQUAL:
			self.addProgramInstructions(EQ)
		case token.NOT_EQUAL:
			self.addProgramInstructions(NE)
		case token.LESS:
			self.addProgramInstructions(LT)
		case token.LESS_OR_EQUAL:
			self.addProgramInstructions(LE)
		case token.GREATER:
			self.addProgramInstructions(GT)
		case token.GREATER_OR_EQUAL:
			self.addProgramInstructions(GE)
		default:
			self.throwSyntaxError(expr.offset, "Unknown operator: %s", expr.operator.String())
		}

	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledAssignExpression(expr *CompiledAssignExpression, putOnStack bool) {
	switch expr.operator {
	case token.ASSIGN:
		self.handlingSetterExpression(expr.left, expr.right, putOnStack)
	case token.ADDITION:
		self.handlingUnaryExpression(expr.left, func() {
			self.handlingGetterExpression(expr.right, true)
			self.addProgramInstructions(Add)
		}, false, putOnStack)
	case token.SUBTRACT:
		self.handlingUnaryExpression(expr.left, func() {
			self.handlingGetterExpression(expr.right, true)
			self.addProgramInstructions(Sub)
		}, false, putOnStack)
	case token.MULTIPLY:
		self.handlingUnaryExpression(expr.left, func() {
			self.handlingGetterExpression(expr.right, true)
			self.addProgramInstructions(Mul)
		}, false, putOnStack)
	case token.DIVIDE:
		self.handlingUnaryExpression(expr.left, func() {
			self.handlingGetterExpression(expr.right, true)
			self.addProgramInstructions(Div)
		}, false, putOnStack)
	case token.REMAINDER:
		self.handlingUnaryExpression(expr.left, func() {
			self.handlingGetterExpression(expr.right, true)
			self.addProgramInstructions(Mod)
		}, false, putOnStack)
	case token.AND_ARITHMETIC:
		self.handlingUnaryExpression(expr.left, func() {
			self.handlingGetterExpression(expr.right, true)
			self.addProgramInstructions(AND)
		}, false, putOnStack)
	case token.OR_ARITHMETIC:
		self.handlingUnaryExpression(expr.left, func() {
			self.handlingGetterExpression(expr.right, true)
			self.addProgramInstructions(OR)
		}, false, putOnStack)
	default:
		self.throwSyntaxError(expr.offset, "Unknown assign operator: %s", expr.operator.String())
	}
}

func (self *Compiler) handlingGetterCompiledFunLiteralExpression(expr *CompiledFunLiteralExpression, putOnStack bool) {
	originProgram := self.program
	self.program = &Program{
		source:       originProgram.source,
		instructions: InstructionArray{},
		sourceMaps:   SourceMapItemArray{{pos: expr.offset}},
	}
	funcScope := self.openScope()
	funcScope.scopeType = ScopeFunction
	funcScope.args = len(expr.parameterList.List)
	enterFunIndex := self.getInstructionSize()
	self.addProgramInstructions(nil)

	if expr.name != nil {
		self.program.functionName = expr.name.Name
	}
	hasInit := false
	for i, binding := range expr.parameterList.List {
		switch target := binding.Target.(type) {
		case *ast.Identifier:
			b, exists := funcScope.bindName(target.Name)
			if exists {
				self.throwSyntaxError(int(target.StartIndex())-1, "Duplicate parameter name not allowed in this context")
			}
			b.isArg = true
			if binding.Initializer == nil {
				continue
			}
			markIndex := self.getInstructionSize()
			self.addProgramInstructions(nil)
			JeqNullIndex := self.getInstructionSize()
			self.addProgramInstructions(nil)
			funcScope.bindings[i].markAccessPointAt(funcScope, markIndex)
			self.setProgramInstruction(markIndex, LoadStackVar(0))
			self.emitVarAssign(target.Name, int(target.StartIndex())-1, self.compileExpression(binding.Initializer))
			self.setProgramInstruction(JeqNullIndex, JeqNull(self.getInstructionSize()-JeqNullIndex))
			hasInit = true
		default:
			self.throwSyntaxError(int(target.StartIndex())-1, "Unsupported BindingElement type: %T", target)
		}
	}

	self.scope.bindName(thisBindingName)

	var enterFunBodyIndex int
	if hasInit {
		self.openScopeNested()
		enterFunBodyIndex = self.getInstructionSize()
		self.addProgramInstructions(nil)
	}

	self.compileDeclarationList(expr.declarationList)
	self.compileStatement(expr.body, false)
	body := expr.body.(*ast.BlockStatement).Body
	lastStatementIndex := len(body) - 1
	var lastStatement ast.Statement
	if lastStatementIndex >= 0 {
		lastStatement = body[lastStatementIndex]
	}
	if _, ok := lastStatement.(*ast.ReturnStatement); !ok {
		self.addProgramInstructions(LoadNull, Ret)
	}

	stackSize, stashSize := funcScope.finaliseVarAlloc(0)
	if stashSize > 0 || funcScope.argsInStash {
		self.setProgramInstruction(enterFunIndex, EnterFunStash{funcScope.argsInStash, stackSize, stashSize, funcScope.args})
	} else {
		self.setProgramInstruction(enterFunIndex, EnterFun{stackSize, funcScope.args})
	}

	if hasInit {
		stackSize, stashSize = 0, 0
		for _, b := range self.scope.bindings {
			if b.inStash {
				stashSize++
			} else {
				stackSize++
			}
		}
		self.setProgramInstruction(enterFunBodyIndex, EnterFunBody{EnterBlock{stackSize, stashSize}})
		self.closeScope()
	}

	self.closeScope()

	newFun := &NewFun{TrimWhitespace(expr.funDefinition), self.program.functionName, self.program}
	self.program = originProgram
	self.addProgramInstructions(newFun)

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledCallExpression(expr *CompiledCallExpression, putOnStack bool) {
	switch callee := expr.callee.(type) {
	case *CompiledIdentifierExpression:
		callee.addSourceMap()
		binding, exists := self.scope.lookupName(callee.name)
		if exists {
			binding.markAccessPoint(self.scope)
			self.addProgramInstructions(LoadStackVar(0))
		} else {
			self.addProgramInstructions(LoadDynamicCallee(callee.name))
		}
	case *CompiledDotExpression:
		self.handlingGetterExpression(callee.left, true)
		self.addProgramInstructions(GetPropCallee(callee.name))
	case *CompiledBracketExpression:
		self.handlingGetterExpression(callee.left, true)
		self.handlingGetterExpression(callee.indexOrName, true)
		self.addProgramInstructions(GetPropOrElemCallee)
	default:
		self.addProgramInstructions(LoadNull)
		self.handlingGetterExpression(callee, true)
	}

	for _, argument := range expr.arguments {
		self.chooseHandlingGetterExpression(argument, true)
	}

	self.addProgramInstructions(Call(len(expr.arguments)))

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledDotExpression(expr *CompiledDotExpression, putOnStack bool) {
	self.handlingGetterExpression(expr.left, true)
	expr.addSourceMap()
	self.addProgramInstructions(GetProp(expr.name))

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledBracketExpression(expr *CompiledBracketExpression, putOnStack bool) {
	self.handlingGetterExpression(expr.left, true)
	self.handlingGetterExpression(expr.indexOrName, true)
	expr.addSourceMap()
	self.addProgramInstructions(GetPropOrElem)

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingSetterExpression(expr CompiledExpression, valueExpr CompiledExpression, putOnStack bool) {
	switch expr := expr.(type) {
	case *CompiledIdentifierExpression:
		self.handlingSetterCompiledIdentifierExpression(expr, valueExpr, putOnStack)
	}
}

func (self *Compiler) handlingSetterCompiledIdentifierExpression(expr *CompiledIdentifierExpression, valueExpr CompiledExpression, putOnStack bool) {
	self.addProgramInstructions(ResolveVar(expr.name))
	self.chooseHandlingGetterExpression(valueExpr, true)
	binding, exists := self.scope.lookupName(expr.name)
	if exists {
		binding.markAccessPoint(self.scope)
		self.addProgramInstructions(PutStackVar(0))
	} else {
		self.addProgramInstructions(PutVar(0))
	}
}

func (self *Compiler) handlingUnaryExpression(expr CompiledExpression, instructionBody func(), postfix bool, putOnStack bool) {
	switch expr := expr.(type) {
	case *CompiledIdentifierExpression:
		self.handlingUnaryCompiledIdentifierExpression(expr, instructionBody, postfix, putOnStack)
	}
}

func (self *Compiler) handlingUnaryCompiledIdentifierExpression(expr *CompiledIdentifierExpression, instructionBody func(), postfix bool, putOnStack bool) {
	self.addProgramInstructions(ResolveVar(expr.name))
	self.chooseHandlingGetterExpression(expr, true)
	instructionBody()
	binding, exists := self.scope.lookupName(expr.name)
	if exists {
		binding.markAccessPoint(self.scope)
		self.addProgramInstructions(PutStackVar(0))
	} else {
		self.addProgramInstructions(PutVar(-1))
	}
}
