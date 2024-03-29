package vm

type Binding struct {
	scope        *Scope
	name         string
	accessPoints map[*Scope]*[]int

	isArg   bool
	inStash bool
}

func (self *Binding) getAccessPointsByScope(scope *Scope) *[]int {
	accessPoints := self.accessPoints[scope]
	if accessPoints != nil {
		return accessPoints
	}
	newAccessPoints := make([]int, 0, 1)
	accessPoints = &newAccessPoints
	self.accessPoints[scope] = accessPoints
	return accessPoints
}

func (self *Binding) markAccessPoint(scope *Scope) {
	self.markAccessPointAt(scope, scope.program.getInstructionSize())
}

func (self *Binding) markAccessPointAt(scope *Scope, pos int) {
	accessPoints := self.getAccessPointsByScope(scope)
	*accessPoints = append(*accessPoints, pos-scope.base)
}

func (self *Binding) moveToStash() {
	if self.isArg && !self.scope.argsInStash {
		self.scope.moveArgsToStash()
	} else {
		self.inStash = true
		self.scope.needStash = true
	}
}

type ScopeType int

const (
	ScopeBlock ScopeType = iota
	ScopeFunction
)

type Scope struct {
	scopeType      ScopeType
	outer          *Scope
	nested         []*Scope
	program        *Program
	bindingMapping map[string]*Binding
	bindings       []*Binding

	base        int
	args        int
	argsInStash bool
	needStash   bool
	isDynamic   bool
}

func (self *Scope) bindName(name string) (*Binding, bool) {
	b, exists := self.bindingMapping[name]
	if exists {
		return b, true
	}
	binding := &Binding{
		self,
		name,
		make(map[*Scope]*[]int),
		false,
		false,
	}
	self.bindings = append(self.bindings, binding)
	self.bindingMapping[name] = binding
	return binding, false
}

func (self *Scope) getBinding(name string) *Binding {
	return self.bindingMapping[name]
}

func (self *Scope) moveArgsToStash() {
	for _, binding := range self.bindings {
		if !binding.isArg {
			break
		}
		binding.inStash = true
	}
	self.argsInStash = true
	self.needStash = true
}

func (self *Scope) lookupName(name string) (*Binding, bool) {
	toStash := false
	for scope := self; scope != nil; scope = scope.outer {
		if binding, exists := scope.bindingMapping[name]; exists && scope.outer != nil {
			if toStash && !binding.inStash {
				binding.moveToStash()
			}
			return binding, true
		}
		if scope.scopeType == ScopeFunction {
			toStash = true
		}
	}
	return nil, false
}

func (self *Scope) nearestFunctionScope() *Scope {
	for s := self; s != nil; s = s.outer {
		if s.scopeType == ScopeFunction {
			return s
		}
	}
	return nil
}

func (self *Scope) needStashDeepLevel(scope *Scope) int {
	deepLevel := 0
	for s := self; s != nil && s != scope; s = s.outer {
		if s.needStash {
			deepLevel++
		}
	}
	return deepLevel
}

func (self *Scope) finaliseVarAlloc(stackOffset int) (int, int) {
	stackIndex, stashIndex := 0, 0
	for i, binding := range self.bindings {
		isThis := binding.name == thisBindingName
		if self.isDynamic || binding.inStash {
			for scope, aps := range binding.accessPoints {
				deepLevel := scope.needStashDeepLevel(self)
				index := (deepLevel << 24) | stackIndex
				program := scope.program
				if isThis {

				} else {
					for _, pc := range *aps {
						pc = pc + scope.base
						instruction := program.getInstruction(pc)
						switch instruction.(type) {
						case InitStackVar:
							program.setProgramInstruction(pc, InitStashVar(index))
						case LoadStackVar:
							program.setProgramInstruction(pc, LoadStashVar(index))
						case PutStackVar:
							program.setProgramInstruction(pc, PutStashVar(index))
						}
					}
				}
			}
			stashIndex++
		} else {
			argsInStash := false
			if scope := self.nearestFunctionScope(); scope != nil {
				argsInStash = scope.argsInStash
			}
			var index int
			if !isThis {
				if i < self.args {
					index = -(i + 1)
				} else {
					stackIndex++
					index = stackIndex + stackOffset
				}
			}
			if isThis {

			} else if argsInStash {
				for scope, aps := range binding.accessPoints {
					program := scope.program
					for _, pc := range *aps {
						pc = pc + scope.base
						instruction := program.getInstruction(pc)
						switch instruction.(type) {
						case InitStackVar:
							program.setProgramInstruction(pc, InitStackVar1(index))
						case LoadStackVar:
							program.setProgramInstruction(pc, LoadStackVar1(index))
						case PutStackVar:
							program.setProgramInstruction(pc, PutStackVar1(index))
						}
					}
				}
			} else {
				for scope, aps := range binding.accessPoints {
					program := scope.program
					for _, pc := range *aps {
						pc = pc + scope.base
						instruction := program.getInstruction(pc)
						switch instruction.(type) {
						case InitStackVar:
							program.setProgramInstruction(pc, InitStackVar(index))
						case LoadStackVar:
							program.setProgramInstruction(pc, LoadStackVar(index))
						case PutStackVar:
							program.setProgramInstruction(pc, PutStackVar(index))
						}
					}
				}
			}
		}
	}
	for _, scope := range self.nested {
		scope.finaliseVarAlloc(stackIndex + stackOffset)
	}
	return stackIndex, stashIndex
}

type EnvRegistry struct {
	fieldNames, methodNames []string
}
type PrivateName struct {
	index    int
	isStatic bool
	isMethod bool
}
type ClassScope struct {
	outer                                *ClassScope
	privateNames                         map[string]*PrivateName
	privateInstanceEnv, privateStaticEnv *EnvRegistry
}
