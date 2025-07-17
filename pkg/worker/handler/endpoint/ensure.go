package endpoint

import (
	"net/http"
	"strconv"

	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	for k, v := range mapping {
		for _, x := range v[h.env.Environment] {
			var sta int
			{
				sta = musSta(x)
			}

			var hlt float64
			if sta == http.StatusOK {
				hlt = 1
			} else {
				h.log.Log(
					"level", "info",
					"message", "observed non-ok status code",
					"url", x,
					"code", strconv.Itoa(sta),
				)
			}

			{
				err = h.reg.Gauge(Metric, hlt, map[string]string{"service": k})
				if err != nil {
					return tracer.Mask(err)
				}
			}
		}
	}

	return nil
}

func musSta(url string) int {
	var err error

	var res *http.Response
	{
		res, err = http.Get(url)
		if err != nil {
			return 0
		}
	}

	{
		defer res.Body.Close()
	}

	return res.StatusCode
}
