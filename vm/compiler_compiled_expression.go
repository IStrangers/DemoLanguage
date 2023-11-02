package vm

import "DemoLanguage/token"

func (self *Compiler) evalConstExpr(expr CompiledExpression) (Value, *Exception) {
	if expr, ok := expr.(CompiledLiteralExpression); ok {
		return expr.value, nil
	}
	evalVM := self.evalVM
	originProgram := self.program
	isNewProgram := false
	if evalVM.program == nil {
		evalVM.program = &Program{
			file: self.program.file,
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

func (self *Compiler) emitLoadValue(value Value, putOnStack bool) {
	if !putOnStack {
		return
	}
	index := self.addProgramValue(value)
	self.addProgramInstructions(LoadVal(index))
}

func (self *Compiler) handlingConstExpression(expr CompiledExpression, putOnStack bool) {
	value, ex := self.evalConstExpr(expr)
	if ex != nil {

	} else {
		self.emitLoadValue(value, putOnStack)
	}
}

func (self *Compiler) handlingGetterExpression(expr CompiledExpression, putOnStack bool) {
	if expr.isConstExpression() {
		self.handlingConstExpression(expr, putOnStack)
		return
	}

	switch expr := expr.(type) {
	case *CompiledLiteralExpression:
		self.handlingGetterCompiledLiteralExpression(expr, putOnStack)
	case *CompiledUnaryExpression:
		self.handlingGetterCompiledUnaryExpression(expr, putOnStack)
	case *CompiledBinaryExpression:
		self.handlingGetterCompiledBinaryExpression(expr, putOnStack)
	}
}

func (self *Compiler) handlingGetterCompiledLiteralExpression(expr *CompiledLiteralExpression, putOnStack bool) {
	self.emitLoadValue(expr.value, putOnStack)
}

func (self *Compiler) handlingGetterCompiledUnaryExpression(expr *CompiledUnaryExpression, putOnStack bool) {
	self.handlingGetterExpression(expr.operand, true)

	switch expr.operator {
	case token.NOT:
		self.addProgramInstructions(Not)
	case token.INCREMENT:
		self.addProgramInstructions(Inc)
	case token.DECREMENT:
		self.addProgramInstructions(Dec)
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingGetterCompiledBinaryExpression(expr *CompiledBinaryExpression, putOnStack bool) {
	self.handlingGetterExpression(expr.left, true)
	self.handlingGetterExpression(expr.right, true)

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
		self.errorAssert(false, expr.offset, "Unknown operator: %s", expr.operator.String())
	}

	if !putOnStack {
		self.addProgramInstructions(Pop)
	}
}

func (self *Compiler) handlingSetterExpression(expr CompiledExpression, putOnStack bool) {

}