package interpreter

type Runtime struct {
	scope *Scope
}

func createRunTime() *Runtime {
	return &Runtime{}
}

func (self *Runtime) openScope() {
	self.scope = &Scope{
		outer: self.scope,
	}
	var stashOuter *Stash
	if self.scope.outer != nil {
		stashOuter = self.scope.outer.stash
	}
	self.scope.stash = &Stash{
		outer:        stashOuter,
		valueMapping: make(map[string]Value),
	}
}

func (self *Runtime) closeScope() {
	self.scope = self.scope.outer
}
