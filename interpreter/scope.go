package interpreter

type Scope struct {
	runtime *Runtime
	outer   *Scope
	stash   *Stash
	this    Objectd
	callee  string
	depth   int
}
