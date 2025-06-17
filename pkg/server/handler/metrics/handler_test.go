package metrics

import (
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/spectagocode/pkg/metrics"
	"github.com/xh3b4sd/logger"
)

func tesHan() metrics.API {
	return New(Config{
		Env: envvar.Env{
			Environment: "foo",
		},
		Log: logger.Fake(),
		Met: recorder.NewMeter(recorder.MeterConfig{
			Env: "testing",
		}),
	})
}
