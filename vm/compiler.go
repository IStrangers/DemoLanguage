package vm

import (
	"DemoLanguage/ast"
	"DemoLanguage/file"
	"fmt"
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
	evalVM  *VM
}

func (self *Compiler) compile(in *ast.Program) {
	self.openScope()
	body := in.Body
	remainingStatements := self.definingUpgrading(body)
	self.compileStatements(remainingStatements, true)
}

func (self *Compiler) definingUpgrading(body []ast.Statement) (remainingStatements []ast.Statement) {
	var funs []*ast.FunStatement
	var funNames []string
	var vars []*ast.VarStatement
	var varNames []string
	for _, statement := range body {
		switch st := statement.(type) {
		case *ast.FunStatement:
			funs = append(funs, st)
			funNames = append(funNames, st.FunLiteral.Name.Name)
		case *ast.VarStatement:
			vars = append(vars, st)
			for _, binding := range st.List {
				varNames = append(varNames, binding.Target.(*ast.Identifier).Name)
			}
		default:
			remainingStatements = append(remainingStatements, st)
		}
	}
	if len(funNames) > 0 || len(varNames) > 0 {
		for _, name := range funNames {
			self.scope.bindName(name)
		}
		for _, name := range varNames {
			self.scope.bindName(name)
		}
		self.addProgramInstructions(&BindDefining{
			funNames,
			varNames,
		})
	}
	self.functionUpgrading(funs)
	self.varUpgrading(vars)
	return
}

func (self *Compiler) functionUpgrading(funs []*ast.FunStatement) {
	for _, fun := range funs {
		self.compileFunStatement(fun)
	}
}

func (self *Compiler) varUpgrading(vars []*ast.VarStatement) {
	for _, v := range vars {
		self.compileVarStatement(v)
	}
}

func (self *Compiler) compileDeclarationList(declarationList []*ast.VariableDeclaration) {
	for _, declaration := range declarationList {
		for _, binding := range declaration.List {
			target := binding.Target
			self.scope.bindName(target.(*ast.Identifier).Name)
		}
	}
}

func (self *Compiler) addProgramValue(value Value) int {
	return self.program.addValue(value)
}

func (self *Compiler) addProgramInstructions(instructions ...Instruction) {
	self.program.addInstructions(instructions...)
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

func (self *Compiler) openBlockScope() {
	self.openScope()
	if outer := self.scope.outer; outer != nil {
		outer.nested = append(outer.nested, self.scope)
	}
	self.scope.base = self.program.getInstructionSize()
}

func (self *Compiler) closeScope() {
	self.scope = self.scope.outer
}

func (self *Compiler) enterVirtualMode() func() {
	originProgram := self.program
	self.program = &Program{
		source: self.program.source,
	}
	self.openScope()
	return func() {
		self.program = originProgram
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
