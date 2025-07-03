package keypair

import (
	"time"
)

func (h *Handler) Cooler() time.Duration {
	return 5 * time.Second
}
