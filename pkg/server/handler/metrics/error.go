package metrics

import (
	"github.com/xh3b4sd/tracer"
)

var requestInvalidError = &tracer.Error{
	Description: "The caller provided an invalid request with its RPC that could not be processed.",
}

var actionInvalidError = &tracer.Error{
	Description: "The caller provided an invalid action with its request that could not be processed.",
}
