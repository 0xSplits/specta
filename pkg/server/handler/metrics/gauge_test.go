package metrics

import (
	"context"
	"testing"

	"github.com/0xSplits/spectagocode/pkg/metrics"
	fuzz "github.com/google/gofuzz"
)

func Test_Server_Handler_Metrics_Gauge_Fuzz(t *testing.T) {
	var han metrics.API
	{
		han = tesHan()
	}

	var fuz *fuzz.Fuzzer
	{
		fuz = fuzz.New()
	}

	for range 1000 {
		var inp *metrics.GaugeI
		{
			inp = &metrics.GaugeI{}
		}

		{
			fuz.Fuzz(inp)
		}

		{
			_, _ = han.Gauge(context.Background(), inp)
		}
	}
}
