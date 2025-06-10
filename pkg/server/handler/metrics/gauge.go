package metrics

import (
	"context"
	"fmt"

	"github.com/0xSplits/spectagocode/pkg/metrics"
)

func (h *Handler) Gauge(ctx context.Context, req *metrics.GaugeI) (*metrics.GaugeO, error) {
	fmt.Printf("/metrics.API/Gauge not implemented\n")
	return &metrics.GaugeO{}, nil
}
