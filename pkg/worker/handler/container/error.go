package container

import (
	"github.com/xh3b4sd/tracer"
)

var invalidAmazonResourceNameError = &tracer.Error{
	Kind: "invalidAmazonResourceNameError",
	Desc: "The exporter expected the ARN format to be [arn:aws:ecs:<region>:<account>:service/<cluster>/<service>], but a different format was found.",
}
