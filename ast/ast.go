package ast

// Node is the base node interface
type Node interface {
}

// Expression node should be implement this interface
type Expression interface {
	Node

	// Dummy methods to identify their types
	expressionNode()
}
