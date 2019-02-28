package parser

import (
	"strings"
)

type errorStrings []string

func (e errorStrings) String() string {
	return strings.Join(e, "\n")
}
