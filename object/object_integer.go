package object

import (
	"fmt"
)

type Integer struct {
	Value int64
}

func (n *Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", int64(i.Value))
}

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type : i.Type(),
		Value: uint64(i.Value),
	}
}
