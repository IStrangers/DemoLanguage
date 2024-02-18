package vm

import (
	"fmt"
	"github.com/istrangers/demolanguage/ast"
	"github.com/istrangers/demolanguage/file"
	"regexp"
)

type CompilerError struct {
	Message string
	File    *file.File
	Offset  int
}

type CompilerSyntaxError struct {
	CompilerError
}

func (self CompilerSyntaxError) Error() string {
	if self.File != nil {
		return fmt.Sprintf("SyntaxError: %s at %s", self.Message, self.File.Position(self.Offset))
	}
	return fmt.Sprintf("SyntaxError: %s", self.Message)
}

type CompilerReferenceError struct {
	CompilerError
}

type Compiler struct {
	program *Program
	scope   *Scope
	block   *Block

	classScope *ClassScope
	evalVM     *VM
}

func CreateCompiler() *Compiler {
	evalVM := CreateVM()
	compiler := &Compiler{
		program: &Program{},
		evalVM:  evalVM,
	}
	return compiler
}

func (self *Compiler) compile(in *ast.Program) {
	self.program.source = in.File

	scope := self.openScope()
	scope.isDynamic = true
	body := in.Body
	declarationList := in.DeclarationList
	remainingStatements := self.definingUpgrading(body, declarationList)
	self.compileStatements(remainingStatements, true)

	scope.finaliseVarAlloc(0)
}

func (self *Compiler) definingUpgrading(body []ast.Statement, declarationList []*ast.VariableDeclaration) (remainingStatements []ast.Statement) {
	var funs []*ast.FunStatement
	var funNames []string
	for _, statement := range body {
		switch st := statement.(type) {
		case *ast.FunStatement:
			funs = append(funs, st)
			funNames = append(funNames, st.FunLiteral.Name.Name)
		default:
			remainingStatements = append(remainingStatements, st)
		}
	}
	if len(funNames) > 0 {
		for _, name := range funNames {
			self.scope.bindName(name)
		}
	}
	self.functionUpgrading(funs)
	varNames := self.compileDeclarationList(declarationList)
	self.addProgramInstructions(&BindDefining{
		funNames,
		varNames,
	})
	return
}

func (self *Compiler) functionUpgrading(funs []*ast.FunStatement) {
	for _, fun := range funs {
		self.compileFunStatement(fun)
	}
}

func (self *Compiler) compileDeclarationList(declarationList []*ast.VariableDeclaration) []string {
	return self.compileScopeDeclarationList(self.scope, declarationList)
}

func (self *Compiler) compileScopeDeclarationList(scope *Scope, declarationList []*ast.VariableDeclaration) []string {
	var varNames []string
	for _, declaration := range declarationList {
		for _, binding := range declaration.List {
			switch t := binding.Target.(type) {
			case *ast.Identifier:
				self.checkScopeVarConflict(scope, t.Name, int(t.StartIndex()-1))
				varNames = append(varNames, t.Name)
			}
		}
	}
	return varNames
}

func (self *Compiler) isScopeDeclared(body []ast.Statement) bool {
	for _, st := range body {
		if _, ok := st.(*ast.VarStatement); ok {
			return true
		}
	}
	return false
}

func (self *Compiler) compileScopeDeclared(body []ast.Statement) bool {
	scopeDeclared := false
	var declarationList []*ast.VariableDeclaration
	for _, st := range body {
		if varSt, ok := st.(*ast.VarStatement); ok {
			declarationList = append(declarationList, &ast.VariableDeclaration{Var: varSt.Var, List: varSt.List})
			scopeDeclared = true
		}
	}
	self.compileDeclarationList(declarationList)
	return scopeDeclared
}

func (self *Compiler) addProgramValue(value Value) int {
	return self.program.addValue(value)
}

func (self *Compiler) addProgramInstructions(instructions ...Instruction) {
	self.program.addInstructions(instructions...)
}

func (self *Compiler) getInstructionSize() int {
	return self.program.getInstructionSize()
}

func (self *Compiler) setProgramInstruction(index int, instruction Instruction) {
	self.program.setProgramInstruction(index, instruction)
}

func (self *Compiler) openScope() *Scope {
	self.scope = &Scope{
		outer:          self.scope,
		program:        self.program,
		bindingMapping: make(map[string]*Binding),
	}
	return self.scope
}

func (self *Compiler) openScopeNested() {
	self.openScope()
	if outer := self.scope.outer; outer != nil {
		outer.nested = append(outer.nested, self.scope)
	}
	self.scope.base = self.getInstructionSize()
}

func (self *Compiler) closeScope() {
	self.scope = self.scope.outer
}

func (self *Compiler) openBlockScope() *Block {
	return self.openBlock(BlockScope)
}

func (self *Compiler) openBlockLoop() *Block {
	return self.openBlock(BlockLoop)
}

func (self *Compiler) openBlockSwitch() *Block {
	return self.openBlock(BlockSwitch)
}

func (self *Compiler) openBlockTry() *Block {
	return self.openBlock(BlockTry)
}

func (self *Compiler) openBlock(blockType BlockType) *Block {
	self.block = &Block{
		outer:     self.block,
		blockType: blockType,
	}
	return self.block
}

func (self *Compiler) closeBlock() {
	index := self.getInstructionSize()
	for _, i := range self.block.breaks {
		self.setProgramInstruction(i, Jump(index-i))
	}
	for _, i := range self.block.continues {
		self.setProgramInstruction(i, Jump(self.block.continueBase-i))
	}
	self.block = self.block.outer
}

func (self *Compiler) findBlockByType(blockTypes []BlockType, isBreak bool) *Block {
	for block := self.block; block != nil; block = self.block.outer {
		for _, blockType := range blockTypes {
			if block.blockType == blockType && (blockType != BlockSwitch || isBreak) {
				return block
			}
		}
	}
	return nil
}

func (self *Compiler) updateEnterBlock(enterBlock *EnterBlock) {
	stackSize, stashSize := 0, 0
	for _, b := range self.scope.bindings {
		if b.inStash {
			stashSize++
		} else {
			stackSize++
		}
	}
	enterBlock.stackSize, enterBlock.stashSize = stackSize, stashSize
}

func (self *Compiler) leaveBlockScope(enterBlock *EnterBlock) {
	self.updateEnterBlock(enterBlock)
	leaveBlock := &LeaveBlock{
		stackSize: enterBlock.stackSize,
		popStash:  enterBlock.stashSize > 0,
	}
	self.addProgramInstructions(leaveBlock)
	for _, pc := range self.block.breaks {
		self.setProgramInstruction(pc, leaveBlock)
	}
	self.block.breaks = nil
	self.closeBlock()
}

func (self *Compiler) openClassScope() *ClassScope {
	self.classScope = &ClassScope{
		outer: self.classScope,
	}
	return self.classScope
}

func (self *Compiler) closeClassScope() {
	self.classScope = self.classScope.outer
}

func (self *Compiler) enterVirtualMode() func() {
	originBlock, originProgram := self.block, self.program
	if originBlock != nil {
		self.block = &Block{
			outer:     originBlock.outer,
			blockType: originBlock.blockType,
		}
	}
	self.program = &Program{
		source: self.program.source,
	}
	self.openScope()
	return func() {
		self.block, self.program = originBlock, originProgram
		self.closeScope()
	}
}

func (self *Compiler) throwSyntaxError(offset int, format string, args ...any) CompiledExpression {
	panic(&CompilerSyntaxError{
		CompilerError{
			File:    self.program.source,
			Offset:  offset,
			Message: fmt.Sprintf(format, args...),
		},
	})
	return nil
}

func (self *Compiler) checkVarConflict(name string, pos int) {
	self.checkScopeVarConflict(self.scope, name, pos)
}

func (self *Compiler) checkScopeVarConflict(scope *Scope, name string, pos int) {
	if _, exists := scope.bindName(name); exists {
		self.throwSyntaxError(pos, "Identifier '%s' has already been declared", name)
	}
}

func TrimWhitespace(content string) string {
	pattern := regexp.MustCompile(`\s+`)
	return pattern.ReplaceAllString(content, " ")
}
