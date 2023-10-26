package interpreter

type Runtime struct {
	global Objectd
	scope  *Scope
}

func createRunTime() *Runtime {
	return &Runtime{
		global: BuiltinGlobalObject(),
	}
}

func (self *Runtime) openScope() {
	self.scope = &Scope{
		runtime: self,
		outer:   self.scope,
	}
	var stashOuter *Stash
	if self.scope.outer != nil {
		stashOuter = self.scope.outer.stash
	}
	self.scope.stash = &Stash{
		runtime:      self,
		outer:        stashOuter,
		valueMapping: make(map[string]Value),
	}
}

func (self *Runtime) closeScope() {
	self.scope = self.scope.outer
}

func (self *Runtime) getStash() *Stash {
	stash := self.scope.stash
	return stash
}
