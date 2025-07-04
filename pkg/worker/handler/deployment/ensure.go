package deployment

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

	var pip []pipeline
	{
		pip, err = h.pipeline(det)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for _, x := range pip {
		lab := map[string]string{
			"platform": "codepipeline",
			"success":  x.suc,
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
			"message", "instrumented pipeline execution",
			"execution", x.eid,
			"latency", x.lat.String(),
			"success", x.suc,
		)
	}

	return nil
}
