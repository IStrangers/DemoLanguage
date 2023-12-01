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
	savedPC := self.program.getInstructionSize()
	evalVM.pc = savedPC
	self.handlingGetterExpression(expr, true)
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
	case *CompiledIdentifierExpression:
		self.handlingGetterCompiledIdentifierExpression(expr, putOnStack)
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

func (self *Compiler) handlingGetterCompiledUnaryExpression(expr *CompiledUnaryExpression, putOnStack bool) {
	switch expr.operator {
	case token.NOT:
		self.chooseHandlingGetterExpression(expr.operand, putOnStack)
		self.addProgramInstructions(Not)
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

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledAssignExpression(expr *CompiledAssignExpression, putOnStack bool) {
	switch expr.operator {
	case token.ASSIGN:
		self.handlingSetterExpression(expr.left, expr.right, putOnStack)
	default:
		self.throwSyntaxError(expr.offset, "Unknown assign operator: %s", expr.operator.String())
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
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
	funcScope.args = len(expr.parameterList.List)
	enterFunIndex := self.program.getInstructionSize()
	self.addProgramInstructions(nil)

	if expr.name != nil {
		self.program.functionName = expr.name.Name
	}
	hasInit := false
	for _, binding := range expr.parameterList.List {
		switch target := binding.Target.(type) {
		case *ast.Identifier:
			_, exists := funcScope.bindName(target.Name)
			if exists {
				self.throwSyntaxError(int(target.StartIndex())-1, "Duplicate parameter name not allowed in this context")
			}
			if binding.Initializer == nil {
				continue
			}
			markIndex := self.program.getInstructionSize()
			self.addProgramInstructions(nil)
			JeqNullIndex := self.program.getInstructionSize()
			self.addProgramInstructions(nil)
			funcScope.getBinding(target.Name).markAccessPointAt(funcScope, markIndex)
			self.setProgramInstruction(markIndex, LoadStackVar(0))
			self.emitVarAssign(target.Name, int(target.StartIndex())-1, self.compileExpression(binding.Initializer))
			self.setProgramInstruction(JeqNullIndex, JeqNull(self.program.getInstructionSize()-JeqNullIndex))
			hasInit = true
		default:
			self.throwSyntaxError(int(target.StartIndex())-1, "Unsupported BindingElement type: %T", target)
		}
	}

	var enterFunBodyIndex int
	if hasInit {
		self.openBlockScope()
		enterFunBodyIndex = self.program.getInstructionSize()
		self.addProgramInstructions(nil)
	}

	self.compileDeclarationList(expr.declarationList)
	self.compileStatement(expr.body, false)
	body := expr.body.(*ast.BlockStatement).Body
	if _, ok := body[len(body)-1].(*ast.ReturnStatement); !ok {
		self.addProgramInstructions(LoadNull, Ret)
	}
	stackSize, stashSize := funcScope.finaliseVarAlloc(0)
	self.setProgramInstruction(enterFunIndex, EnterFun{stackSize, funcScope.args})

	if hasInit {
		self.setProgramInstruction(enterFunBodyIndex, EnterFunBody{EnterBlock{stackSize, stashSize}})
		self.closeScope()
	}

	self.closeScope()

	newFun := &NewFun{expr.funDefinition, self.program.functionName, self.program}
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
		self.addProgramInstructions(LoadDynamicCallee(callee.name))
	}

	for _, argument := range expr.arguments {
		self.chooseHandlingGetterExpression(argument, true)
	}

	self.addProgramInstructions(Call(len(expr.arguments)))

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
	self.chooseHandlingGetterExpression(valueExpr, putOnStack)
	self.addProgramInstructions(PutVar(0))
}

func (self *Compiler) handlingUnaryExpression(expr CompiledExpression, instructionBody func(), postfix bool, putOnStack bool) {
	switch expr := expr.(type) {
	case *CompiledIdentifierExpression:
		self.handlingUnaryCompiledIdentifierExpression(expr, instructionBody, postfix, putOnStack)
	}
}

func (self *Compiler) handlingUnaryCompiledIdentifierExpression(expr *CompiledIdentifierExpression, instructionBody func(), postfix bool, putOnStack bool) {
	self.addProgramInstructions(ResolveVar(expr.name))
	self.chooseHandlingGetterExpression(expr, putOnStack)
	instructionBody()
	self.addProgramInstructions(PutVar(-1))
}
