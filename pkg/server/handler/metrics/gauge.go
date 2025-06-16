package metrics

import (
	"context"

	"github.com/0xSplits/spectagocode/pkg/metrics"
)

func (h *Handler) Gauge(ctx context.Context, req *metrics.GaugeI) (*metrics.GaugeO, error) {
	h.log.Log(
		"level", "info",
		"message", "endpoint not implemented",
		"endpoint", "/metrics.API/Gauge",
	)

	return &metrics.GaugeO{}, nil
}
