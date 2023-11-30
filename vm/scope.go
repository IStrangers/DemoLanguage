package vm

type Binding struct {
	scope *Scope
	name  string
}

type Scope struct {
	outer          *Scope
	program        *Program
	bindingMapping map[string]*Binding
}

func (self *Scope) bindName(name string) (*Binding, bool) {
	_, exists := self.bindingMapping[name]
	binding := &Binding{
		self,
		name,
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
