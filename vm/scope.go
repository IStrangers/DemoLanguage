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
