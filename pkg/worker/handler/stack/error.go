package stack

import (
	"github.com/xh3b4sd/tracer"
)

var missingRootStackError = &tracer.Error{
	Kind: "missingRootStackError",
	Desc: "The exporter expected to find exactly one root stack, but no stack was found for teh given environment.",
}
