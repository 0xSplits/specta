package build

import (
	"time"
)

func (h *Handler) Cooler() time.Duration {
	return 1 * time.Minute
}
