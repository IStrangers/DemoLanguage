package interpreter

type ClassObject interface {
	ofLiteral(object Objectd) string
	getValue(object Objectd, property Value, args ...Value) Value
	setValue(object Objectd, property Value, values ...Value)
}
