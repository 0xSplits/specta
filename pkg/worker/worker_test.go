package worker

import (
	"slices"
	"testing"
	"time"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/specta/pkg/runtime"
	"github.com/0xSplits/specta/pkg/worker/handler"
	"github.com/xh3b4sd/logger"
)

func Test_Worker_handler_execution(t *testing.T) {
	var in1 chan string
	var in2 chan string
	var out chan string
	{
		in1 = make(chan string, 1)
		in2 = make(chan string, 3)
		out = make(chan string, 6)
	}

	var wrk *Worker
	{
		wrk = New(Config{
			Env: envvar.Env{Environment: "testing"},
			Han: []handler.Interface{
				&testHandler{inp: in1, out: out, coo: time.Hour},
				&testHandler{inp: in2, out: out, coo: 0},
			},
			Log: logger.Fake(),
			Met: recorder.NewMeter(recorder.MeterConfig{
				Env: "testing",
				Ver: runtime.Tag(),
			}),
		})
	}

	{
		go wrk.Daemon()
	}

	{
		time.Sleep(time.Millisecond)
	}

	go func() {
		in1 <- "one"
		in1 <- "one"
		in1 <- "one"
	}()

	go func() {
		in2 <- "two"
		in2 <- "two"
		in2 <- "two"
	}()

	var exp []string
	{
		exp = []string{
			"one",

			"two",
			"two",
			"two",
		}
	}

	var act []string
	for x := range out {
		{
			act = append(act, x)
		}

		if len(act) == len(exp) {
			close(out)
		}
	}

	select {
	case <-out:
		{
			slices.Sort(act)
			slices.Sort(exp)
		}

		if !slices.Equal(act, exp) {
			t.Fatalf("expected %#v got %#v", exp, act)
		}
	case <-time.After(time.Second):
		t.Fatal("test timeout")
	}
}

type testHandler struct {
	coo time.Duration
	inp chan string
	out chan string
}

func (h *testHandler) Cooler() time.Duration {
	return h.coo
}

func (h *testHandler) Ensure() error {
	var s string
	{
		s = <-h.inp
	}

	{
		h.out <- s
	}

	return nil
}
