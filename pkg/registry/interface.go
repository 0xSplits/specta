package registry

import "go.opentelemetry.io/otel/metric"

type Interface interface {
	Counter(string) (metric.Float64Counter, map[string]struct{}, bool)
	Gauge(string) (metric.Float64Gauge, map[string]struct{}, bool)
	Histogram(string) (metric.Float64Histogram, map[string]struct{}, bool)
}
