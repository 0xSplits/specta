package metrics

import (
	"github.com/0xSplits/spectagocode/pkg/metrics"
	"github.com/xh3b4sd/tracer"
)

const (
	// max is the maximum amount of actions any given request is allowed to
	// provide at once.
	max = 100
)

type Request interface {
	GetAction() []*metrics.Action
}

// verify ensures that the given request is properly configured within the
// bounds of the allowed parameters. E.g. requests must provide at least one
// action.
func verify(req Request) error {
	if req == nil {
		return tracer.Maskf(requestInvalidError, "request must not be empty")
	}

	var act []*metrics.Action
	{
		act = req.GetAction()
	}

	{
		if len(act) == 0 {
			return tracer.Maskf(requestInvalidError, "request must have at least one action")
		}
		if len(act) > max {
			return tracer.Maskf(requestInvalidError, "request must have at most %d actions", max)
		}
	}

	for i, x := range act {
		if x == nil {
			return tracer.Maskf(actionInvalidError, "action[%d] must not be empty", i)
		}
		if x.GetMetric() == "" {
			return tracer.Maskf(actionInvalidError, "action[%d] must not have empty metric name", i)
		}
		if len(x.GetMetric()) > 255 {
			return tracer.Maskf(actionInvalidError, "action[%d] must not have metric name longer than 255 characters", i)
		}
		if x.GetNumber() < 0 {
			return tracer.Maskf(actionInvalidError, "action[%d] must not have negative number", i)
		}
	}

	return nil
}
