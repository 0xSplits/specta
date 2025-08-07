package endpoint

import (
	"net/http"
	"slices"

	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

type endpoint struct {
	// hlt is the service endpoint health, either 0.0 or 1.0.
	hlt float64
	// lab is the respective service label, e.g. explorer or specta.
	lab string
}

func (h *Handler) endpoint(det []detail) ([]endpoint, error) {
	var end []endpoint
	var err error

	fnc := func(_ int, d detail) error {
		var lis []float64

		for _, x := range d.url {
			var hlt int
			{
				hlt = musHlt(x)
			}

			if hlt != 1 {
				h.log.Log(
					"level", "info",
					"message", "observed non-ok status code",
					"url", x,
				)
			}

			{
				lis = append(lis, float64(hlt))
			}
		}

		{
			end = append(end, endpoint{
				hlt: slices.Min(lis),
				lab: d.lab,
			})
		}

		return nil
	}

	{
		err = parallel.Slice(det, fnc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return end, nil
}

func musHlt(url string) int {
	var err error

	var res *http.Response
	{
		res, err = http.Get(url)
		if err != nil {
			return 0
		}
	}

	{
		defer res.Body.Close() //nolint:errcheck
	}

	if res.StatusCode == http.StatusOK {
		return 1
	}

	return 0
}
