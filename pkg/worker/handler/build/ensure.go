package build

import (
	"strconv"

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

	var wor []workflow
	{
		wor, err = h.workflow(det)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for _, x := range wor {
		lab := map[string]string{
			"platform":   "github",
			"repository": x.rep,
			"success":    x.suc,
			"workflow":   x.kin,
		}

		err = h.reg.Counter(MetricTotal, 1, lab)
		if err != nil {
			return tracer.Mask(err)
		}

		err = h.reg.Histogram(MetricDuration, x.lat.Seconds(), lab)
		if err != nil {
			return tracer.Mask(err)
		}

		h.log.Log(
			"level", "debug",
			"message", "instrumented build container",
			"latency", x.lat.String(),
			"repository", x.rep,
			"run", strconv.FormatInt(x.rid, 10),
			"success", x.suc,
		)
	}

	return nil
}
