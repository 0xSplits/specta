package metrics

import (
	"context"
	"testing"

	"github.com/0xSplits/spectagocode/pkg/metrics"
	fuzz "github.com/google/gofuzz"
)

func Test_Server_Handler_Metrics_Counter_Fuzz(t *testing.T) {
	var han metrics.API
	{
		han = tesHan()
	}

	var fuz *fuzz.Fuzzer
	{
		fuz = fuzz.New()
	}

	for range 1000 {
		var inp *metrics.CounterI
		{
			inp = &metrics.CounterI{}
		}

		{
			fuz.Fuzz(inp)
		}

		{
			_, _ = han.Counter(context.Background(), inp)
		}
	}
}
