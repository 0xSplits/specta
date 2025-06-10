package metrics

import (
	"context"
	"fmt"

	"github.com/0xSplits/spectagocode/pkg/metrics"
)

func (h *Handler) Histogram(ctx context.Context, req *metrics.HistogramI) (*metrics.HistogramO, error) {
	fmt.Printf("/metrics.API/Histogram not implemented\n")
	return &metrics.HistogramO{}, nil
}
