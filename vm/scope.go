package vm

type Scope struct {
	outer   *Scope
	program *Program
}
