package vm

import (
	"github.com/istrangers/demolanguage/ast"
	"github.com/istrangers/demolanguage/token"
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
		self.handlingGetterCompiledCallExpression(expr, false, putOnStack)
	case *CompiledDotExpression:
		self.handlingGetterCompiledDotExpression(expr, putOnStack)
	case *CompiledBracketExpression:
		self.handlingGetterCompiledBracketExpression(expr, putOnStack)
	case *CompiledClassLiteralExpression:
		self.handlingGetterCompiledClassLiteralExpression(expr, putOnStack)
	case *CompiledNewExpression:
		self.handlingGetterCompiledNewExpression(expr, putOnStack)
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
	funProgram, argNum := self.compileFunProgram(expr)
	newFun := &NewFun{TrimWhitespace(expr.funDefinition), funProgram.functionName, argNum, funProgram}
	self.addProgramInstructions(newFun)

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) compileFunProgram(expr *CompiledFunLiteralExpression) (*Program, int) {
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
			self.setProgramInstruction(markIndex, LoadStackVar(0))

			self.chooseHandlingGetterExpression(self.compileExpression(binding.Initializer), true)
			funcScope.bindings[i].markAccessPointAt(funcScope, markIndex)
			funcScope.bindings[i].markAccessPoint(funcScope)
			self.addProgramInstructions(InitStackVar(0))
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
		self.compileScopeDeclarationList(funcScope, expr.declarationList)
	}

	self.compileDeclarationList(expr.declarationList)
	body := expr.body.Body
	self.compileStatements(body, false)
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
		enterBlock := EnterBlock{}
		self.updateEnterBlock(&enterBlock)
		self.setProgramInstruction(enterFunBodyIndex, EnterFunBody{EnterBlock: enterBlock})
		self.closeScope()
	}

	funProgram := self.program
	self.closeScope()
	self.program = originProgram

	return funProgram, len(expr.parameterList.List)
}

func (self *Compiler) handlingGetterCompiledCallExpression(expr *CompiledCallExpression, isNewCall bool, putOnStack bool) {
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

	if isNewCall {
		self.addProgramInstructions(New(len(expr.arguments)))
	} else {
		self.addProgramInstructions(Call(len(expr.arguments)))
	}

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

func (self *Compiler) handlingGetterCompiledClassLiteralExpression(expr *CompiledClassLiteralExpression, putOnStack bool) {
	self.openScopeNested()

	//enterBlock := &EnterBlock{}
	//markIndex := self.getInstructionSize()
	//self.addProgramInstructions(enterBlock)
	self.block = self.openBlockScope()
	self.scope.bindName(expr.name.Name)
	self.openClassScope()

	newClass := &NewClass{
		name:   expr.name.Name,
		source: expr.classDefinition,
	}

	var newClassInstruction Instruction
	isDerivedClass := false
	if expr.superClass != nil {
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
	for _, declaration := range expr.body {
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
				program, argNum := self.compileConstructor(decl.Body, isDerivedClass)
				newClass.constructors = append(newClass.constructors, &Constructor{
					decl.Body.FunDefinition,
					argNum,
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

	if staticCount > 0 {
		newClass.staticInit = self.compileDeclarations("<static_initializer>", staticBlocks, staticFieldDecls, staticMethodDecls)
	}

	if instanceCount > 0 {
		newClass.instanceInit = self.compileDeclarations("<instance_members_initializer>", nil, instanceFieldDecls, instanceMethodDecls)
	}

	if isDerivedClass {
		self.handlingGetterExpression(self.compileExpression(expr.superClass), true)
	}
	expr.addSourceMap()
	self.addProgramInstructions(newClassInstruction)

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}

	self.closeClassScope()
	self.closeScope()
}

func (self *Compiler) handlingGetterCompiledNewExpression(expr *CompiledNewExpression, putOnStack bool) {
	callExpression := expr.callExpression
	self.handlingGetterCompiledCallExpression(callExpression, true, true)

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingSetterExpression(expr CompiledExpression, valueExpr CompiledExpression, putOnStack bool) {
	switch expr := expr.(type) {
	case *CompiledIdentifierExpression:
		self.handlingSetterCompiledIdentifierExpression(expr, valueExpr, putOnStack)
	case *CompiledDotExpression:
		self.handlingSetterCompiledDotExpression(expr, valueExpr, putOnStack)
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

func (self *Compiler) handlingSetterCompiledDotExpression(expr *CompiledDotExpression, valueExpr CompiledExpression, putOnStack bool) {
	self.handlingGetterExpression(expr.left, true)
	self.handlingGetterExpression(valueExpr, true)
	expr.addSourceMap()
	self.addProgramInstructions(AddProp(expr.name))

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingUnaryExpression(expr CompiledExpression, instructionBody func(), postfix bool, putOnStack bool) {
	switch expr := expr.(type) {
	case *CompiledIdentifierExpression:
		self.handlingUnaryCompiledIdentifierExpression(expr, instructionBody, postfix, putOnStack)
	}
}

func (self *Compiler) handlingUnaryCompiledIdentifierExpression(expr *CompiledIdentifierExpression, instructionBody func(), postfix bool, putOnStack bool) {
	binding, exists := self.scope.lookupName(expr.name)
	if exists {
		self.chooseHandlingGetterExpression(expr, true)
		instructionBody()
		binding.markAccessPoint(self.scope)
		self.addProgramInstructions(PutStackVar(0))
	} else {
		self.addProgramInstructions(ResolveVar(expr.name))
		self.chooseHandlingGetterExpression(expr, true)
		instructionBody()
		self.addProgramInstructions(PutVar(-1))
	}
}
