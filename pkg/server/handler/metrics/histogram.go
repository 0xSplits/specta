package metrics

import (
	"context"

	"github.com/0xSplits/specta/pkg/status"
	"github.com/0xSplits/spectagocode/pkg/metrics"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Histogram(ctx context.Context, req *metrics.HistogramI) (*metrics.HistogramO, error) {
	{
		err := verify(req)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	for _, x := range req.GetAction() {
		err := h.reg.Histogram(x.GetMetric(), x.GetNumber(), x.GetLabels())
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var res *metrics.HistogramO
	{
		res = &metrics.HistogramO{}
	}

	for range req.GetAction() {
		res.Result = append(res.Result, &metrics.Result{
			Status: status.Success,
		})
	}

	return res, nil
}
