package object

import (
	"hash/fnv"
)

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJECT
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Inspect()))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}
