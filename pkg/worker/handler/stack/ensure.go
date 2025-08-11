package stack

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

	var sta []stack
	{
		sta, err = h.stack(det)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for _, x := range sta {
		err = h.reg.Gauge(Metric, x.hlt, map[string]string{"stack": x.lab})
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
