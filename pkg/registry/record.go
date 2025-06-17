package registry

import (
	"slices"

	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/attribute"
)

// record validates and executes the data collection for the given metric, if
// that metric is registered in the provided whitelist of recorder interfaces.
// The correct metric names must be registered for the respective counters,
// gauges and histograms. The provided labels must be registered too, by key and
// value respectively.
func (r *Registry) record(wht map[string]recorder.Interface, nam string, val float64, lab map[string]string) error {
	var rec recorder.Interface
	var exi bool

	{
		rec, exi = wht[nam]
		if !exi {
			return tracer.Maskf(metricNameWhitelistError, "name = %s", nam)
		}
	}

	for k, v := range lab {
		var lis []string
		{
			lis, exi = rec.Labels()[k]
			if !exi {
				return tracer.Maskf(labelKeyWhitelistError, "key = %s", k)
			}
		}

		{
			exi = slices.Contains(lis, v)
			if !exi {
				return tracer.Maskf(labelValueWhitelistError, "value = %s", v)
			}
		}
	}

	var att []attribute.KeyValue
	{
		att = append(att, attribute.String("env", r.env.Environment))
	}

	for k, v := range lab {
		att = append(att, attribute.String(k, v))
	}

	{
		rec.Record(val, att...)
	}

	return nil
}
