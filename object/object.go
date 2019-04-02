package object

//
const (
	RETURN_VALUE_OBJ = "RETURN_VALUE_OBJ"
	NIL_OBJECT       = "NIL_OBJECT"
)

//
type ObjectType string

//
type Object interface {
	Type() ObjectType
	Inspect() string
}
