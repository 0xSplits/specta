package privatekey

import (
	"fmt"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/specta/pkg/registry"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	env envvar.Env
	log logger.Interface
	reg registry.Interface
}

func New(c Config) *Handler {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Met == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Met must not be empty", c)))
	}

	cou := map[string]recorder.Interface{}
	gau := map[string]recorder.Interface{}
	his := map[string]recorder.Interface{}

	{
		// metrics can be registered here
	}

	var reg registry.Interface
	{
		reg = registry.New(registry.Config{
			Env: c.Env,
			Log: c.Log,

			Cou: cou,
			Gau: gau,
			His: his,
		})
	}

	return &Handler{
		env: c.Env,
		log: c.Log,
		reg: reg,
	}
}
