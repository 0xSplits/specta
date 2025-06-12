package registry

import (
	"fmt"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	Log logger.Interface
	Met metric.Meter
}

type Registry struct {
	log logger.Interface
	met metric.Meter
}

func New(c Config) *Registry {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Met == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Met must not be empty", c)))
	}

	// TODO register whitelisted metrics
	// counter, err := c.Met.Float64Counter("foo", metric.WithDescription("a simple counter"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// gauge, err := c.Met.Float64Gauge("foo", metric.WithDescription("a simple gauge"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// historgam, err := c.Met.Float64Histogram("foo", metric.WithDescription("a simple histogram"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return &Registry{
		log: c.Log,
		met: c.Met,
	}
}

func (r *Registry) Counter(string) (metric.Float64Counter, map[string]struct{}, error) {
	// TODO implement getter for counter
	return nil, nil, nil
}

func (r *Registry) Gauge(string) (metric.Float64Gauge, map[string]struct{}, error) {
	// TODO implement getter for gauge
	return nil, nil, nil
}

func (r *Registry) Histogram(string) (metric.Float64Histogram, map[string]struct{}, error) {
	// TODO implement getter for histogram
	return nil, nil, nil
}
