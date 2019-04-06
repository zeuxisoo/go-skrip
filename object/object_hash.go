package object

import (
	"fmt"
	"bytes"
	"strings"
	"sort"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJECT
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	var pairs []string

	for _, pair := range h.Pairs {
		pairs = append(
			pairs,
			fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()),
		)
	}

	sort.Strings(pairs)

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
