package interpreter

type ClassObject interface {
	getValue(object Objectd, property Value, args ...Value) Value
	setValue(object Objectd, property Value, values ...Value)
}
