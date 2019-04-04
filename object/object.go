package object

//
const (
	RETURN_VALUE_OBJECT = "RETURN_VALUE_OBJECT"
	NIL_OBJECT       	= "NIL_OBJECT"
	ERROR_OBJECT		= "ERROR_OBJECT"
	INTEGER_OBJECT		= "INTEGER_OBJECT"
	FLOAT_OBJECT		= "FLOAT_OBJECT"
	STRING_OBJECT		= "STRING_OBJECT"
)

//
type ObjectType string

//
type Object interface {
	Type() ObjectType
	Inspect() string
}
