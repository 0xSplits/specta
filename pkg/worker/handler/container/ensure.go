package container

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

	var ser []service
	{
		ser, err = h.service(det)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for _, x := range ser {
		err = h.reg.Gauge(Metric, x.hlt, map[string]string{"service": x.lab})
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
