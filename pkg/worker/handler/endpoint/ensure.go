package endpoint

import (
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	var det []detail
	{
		det, err = h.detail()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var end []endpoint
	{
		end, err = h.endpoint(det)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for _, x := range end {
		err = h.reg.Gauge(Metric, x.hlt, map[string]string{"service": x.lab})
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
