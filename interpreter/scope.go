package interpreter

type Scope struct {
	outer *Scope
	stash *Stash
}
