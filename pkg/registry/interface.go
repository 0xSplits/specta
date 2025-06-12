package registry

import "go.opentelemetry.io/otel/metric"

// Interface provides a mechanism for the server and worker handlers to get the
// metrics they are interested in. This abstract metrics registry defines the
// whitelist of available metrics. The consumers, here server and worker
// handlers, will use the registry to get their desired metrics on demand. It is
// important for the registry to only provide predefined metrics for timeseries
// that we want to track deliberately. E.g. it should not be possible for the
// outside world to create metrics or labels arbitrarily based on user input.
type Interface interface {
	// Counter performs a registry lookup for the counter metric matching the
	// provided metric name.
	//
	//     inp[0] the name of the counter metric to fetch
	//
	//     out[0] the interface of the counter metric to work with
	//     out[1] the set of whitelisted labels to work with, if any
	//     out[2] the error describing why the lookup failed, if any
	//
	Counter(string) (metric.Float64Counter, map[string]struct{}, error)

	// Gauge performs a registry lookup for the gauge metric matching the provided
	// metric name.
	//
	//     inp[0] the name of the gauge metric to fetch
	//
	//     out[0] the interface of the gauge metric to work with
	//     out[1] the set of whitelisted labels to work with, if any
	//     out[2] the error describing why the lookup failed, if any
	//
	Gauge(string) (metric.Float64Gauge, map[string]struct{}, error)

	// Histogram performs a registry lookup for the histogram metric matching the
	// provided metric name.
	//
	//     inp[0] the name of the histogram metric to fetch
	//
	//     out[0] the interface of the histogram metric to work with
	//     out[1] the set of whitelisted labels to work with, if any
	//     out[2] the error describing why the lookup failed, if any
	//
	Histogram(string) (metric.Float64Histogram, map[string]struct{}, error)
}
