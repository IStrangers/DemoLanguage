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

func (self *Runtime) openScope(this Objectd) {
	self.scope = &Scope{
		runtime: self,
		outer:   self.scope,
		this:    this,
		depth:   0,
	}
	var stashOuter *Stash
	if self.scope.outer != nil {
		self.scope.depth = self.scope.outer.depth + 1
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
