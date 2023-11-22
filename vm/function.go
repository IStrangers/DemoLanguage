package vm

type BaseFunObject struct {
	BaseObject
	funDefinition string
	program       *Program
}

type FunObject struct {
	BaseFunObject
}
