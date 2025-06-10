package failure

import (
	"context"
	"fmt"

	"github.com/twitchtv/twirp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Log logger.Interface
}

type Interceptor struct {
	log logger.Interface
}

func New(c Config) *Interceptor {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Interceptor{
		log: c.Log,
	}
}

func (i *Interceptor) Method(nex twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		res, err := nex(ctx, req)
		if err != nil {
			e, o := err.(*tracer.Error)
			if o {
				i.log.Log(
					"level", "error",
					"message", e.Error(),
					"code", e.Code,
					"description", e.Desc,
					"docs", e.Docs,
					"kind", e.Kind,
					"stack", tracer.Stack(e),
				)

				c := twirp.ErrorCode(e.Code)
				if c == "" {
					c = twirp.Internal
				}

				var t twirp.Error
				{
					t = twirp.NewError(c, e.Error())
				}

				if e.Desc != "" {
					t = t.WithMeta("desc", e.Desc)
				}

				if e.Kind != "" {
					t = t.WithMeta("kind", e.Kind)
				}

				return nil, t
			} else {
				i.log.Log(
					"level", "error",
					"message", err.Error(),
				)

				return nil, err
			}
		}

		return res, nil
	}
}
