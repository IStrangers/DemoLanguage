package vm

type Binding struct {
	scope        *Scope
	name         string
	accessPoints map[*Scope]*[]int
}

func (self Binding) getAccessPointsByScope(scope *Scope) *[]int {
	accessPoints := self.accessPoints[scope]
	if accessPoints != nil {
		return accessPoints
	}
	newAccessPoints := make([]int, 0, 1)
	accessPoints = &newAccessPoints
	self.accessPoints[scope] = accessPoints
	return accessPoints
}

func (self Binding) markAccessPoint(scope *Scope) {
	self.markAccessPointAt(scope, scope.program.getInstructionSize())
}

func (self Binding) markAccessPointAt(scope *Scope, pos int) {
	accessPoints := self.getAccessPointsByScope(scope)
	*accessPoints = append(*accessPoints, pos-scope.base)
}

type Scope struct {
	outer          *Scope
	nested         []*Scope
	program        *Program
	bindingMapping map[string]*Binding
	bindings       []*Binding

	base int
	args int
}

func (self *Scope) bindName(name string) (*Binding, bool) {
	_, exists := self.bindingMapping[name]
	binding := &Binding{
		self,
		name,
		make(map[*Scope]*[]int),
	}
	self.bindings = append(self.bindings, binding)
	self.bindingMapping[name] = binding
	return binding, exists
}

func (self *Scope) lookupName(name string) (*Binding, bool) {
	for scope := self; scope != nil; scope = scope.outer {
		if binding, exists := scope.bindingMapping[name]; exists && scope.outer != nil {
			return binding, true
		}
	}
	return nil, false
}

func (self *Scope) getBinding(name string) *Binding {
	return self.bindingMapping[name]
}

func (self *Scope) finaliseVarAlloc(stackOffset int) (int, int) {
	stackIndex, stashIndex := 0, 0
	for i, binding := range self.bindings {
		var index int
		if binding.name != thisBindingName {
			if i <= self.args {
				index = -(i + 1)
			} else {
				stackIndex++
				index = stackIndex + stackOffset
			}
		}
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
				}
			}
		}
	}
	for _, scope := range self.nested {
		scope.finaliseVarAlloc(stackIndex + stackOffset)
	}
	return stackIndex, stashIndex
}
