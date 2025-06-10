package worker

import (
	"fmt"
	"time"

	"github.com/0xSplits/specta/pkg/worker/handler"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Han are the worker specific handlers implementing the actual business
	// logic.
	Han []handler.Interface
	Log logger.Interface
}

type Worker struct {
	han []handler.Interface
	log logger.Interface
}

func New(c Config) *Worker {
	if len(c.Han) == 0 {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Han must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Worker{
		han: c.Han,
		log: c.Log,
	}
}

func (w *Worker) Daemon() {
	w.log.Log(
		"level", "info",
		"message", "worker is executing tasks",
	)

	for {
		for _, h := range w.han {
			err := h.Ensure()
			if err != nil {
				w.log.Log(
					"level", "error",
					"message", err.Error(),
				)
			}
		}

		{
			time.Sleep(5 * time.Second)
		}
	}
}
