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
	accessPoints := self.getAccessPointsByScope(scope)
	*accessPoints = append(*accessPoints, scope.program.getInstructionSize())
}

type Scope struct {
	outer          *Scope
	program        *Program
	bindingMapping map[string]*Binding
	args           int
}

func (self *Scope) bindName(name string) (*Binding, bool) {
	_, exists := self.bindingMapping[name]
	binding := &Binding{
		self,
		name,
		make(map[*Scope]*[]int),
	}
	self.bindingMapping[name] = binding
	return binding, exists
}

func (self *Scope) lookupName(name string) (*Binding, bool) {
	for scope := self; scope != nil; scope = scope.outer {
		if binding, exists := scope.bindingMapping[name]; exists {
			return binding, true
		}
	}
	return nil, false
}

func (self *Scope) finaliseVarAlloc(stackOffset int) (int, int) {
	stackIndex, stashIndex := 0, 0
	i := 1
	for _, binding := range self.bindingMapping {
		var index int
		if i <= self.args {
			index = -i
		} else {
			stackIndex++
			index = stackIndex + stackOffset
		}
		for scope, aps := range binding.accessPoints {
			program := scope.program
			for _, pc := range *aps {
				instruction := program.getInstruction(pc)
				switch instruction.(type) {
				case InitStackVar:
					program.setProgramInstruction(pc, InitStackVar(index))
				case LoadStackVar:
					program.setProgramInstruction(pc, LoadStackVar(index))
				}
			}
		}
		i++
	}
	return stackIndex, stashIndex
}
