package stack

import (
	"github.com/xh3b4sd/tracer"
)

var missingRootStackError = &tracer.Error{
	Kind: "missingRootStackError",
	Desc: "The exporter expected to find exactly one root stack, but no stack name matched against the internal mapping.",
}

var tooManyRootStacksError = &tracer.Error{
	Kind: "tooManyRootStacksError",
	Desc: "The exporter expected to find exactly one root stack, but more than one stack names matched against the internal mapping.",
}
