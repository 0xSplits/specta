package endpoint

import (
	"net/http"

	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	var res *http.Response
	{
		res, err = http.Get(h.end)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		defer res.Body.Close()
	}

	var hlt float64
	if res.StatusCode == http.StatusOK {
		hlt = 1
	}

	{
		err = h.reg.Gauge(Metric, hlt, map[string]string{"service": h.ser})
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
