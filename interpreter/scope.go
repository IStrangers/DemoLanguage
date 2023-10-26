package interpreter

type Scope struct {
	runtime *Runtime
	outer   *Scope
	stash   *Stash
}
