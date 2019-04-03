package object

import (
	"strconv"
)

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJECT
}

func (f *Float) Inspect() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}
