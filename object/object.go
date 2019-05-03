package object

//
const (
	RETURN_VALUE_OBJECT = "RETURN_VALUE_OBJECT"
	NIL_OBJECT          = "NIL_OBJECT"
	ERROR_OBJECT        = "ERROR_OBJECT"
	INTEGER_OBJECT      = "INTEGER_OBJECT"
	FLOAT_OBJECT        = "FLOAT_OBJECT"
	STRING_OBJECT       = "STRING_OBJECT"
	BUILTIN_OBJECT      = "BUILTIN_OBJECT"
	BOOLEAN_OBJECT      = "BOOLEAN_OBJECT"
	ARRAY_OBJECT        = "ARRAY_OBJECT"
	HASH_OBJECT         = "HASH_OBJECT"
	FUNCTION_OBJECT     = "FUNCTION_OBJECT"
	BREAK_OBJECT        = "BREAK_OBJECT"
)

//
type ObjectType string

//
type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}
