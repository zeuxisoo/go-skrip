package ast

// Node is the base node interface
type Node interface {
	TokenLiteral() string
	String() string
}

// Expression node should be implement this interface
type Expression interface {
	Node
	expressionNode() // Dummy methods to identify their types
}
