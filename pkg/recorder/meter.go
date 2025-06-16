package recorder

import (
	"fmt"

	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type MeterConfig struct {
	Env string
}

func NewMeter(c MeterConfig) otelmetric.Meter {
	if c.Env == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}

	exp, err := prometheus.New()
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exp),
	).Meter(fmt.Sprintf("specta.%s.splits.org", c.Env))
}
