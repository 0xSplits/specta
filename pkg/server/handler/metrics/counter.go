package metrics

import (
	"context"
	"fmt"

	"github.com/0xSplits/spectagocode/pkg/metrics"
)

func (h *Handler) Counter(ctx context.Context, req *metrics.CounterI) (*metrics.CounterO, error) {
	fmt.Printf("/metrics.API/Counter not implemented\n")
	return &metrics.CounterO{}, nil
}
