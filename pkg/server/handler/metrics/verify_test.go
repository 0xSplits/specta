package metrics

import (
	"errors"
	"fmt"
	"testing"

	"github.com/0xSplits/spectagocode/pkg/metrics"
)

func Test_Server_Handler_Metrics_verify_failure(t *testing.T) {
	testCases := []struct {
		req Request
		err error
	}{
		// Case 000, nil request
		{
			req: nil,
			err: requestInvalidError,
		},
		// Case 001, no actions
		{
			req: &metrics.CounterI{},
			err: requestInvalidError,
		},
		// Case 002, too many actions
		{
			req: &metrics.CounterI{
				Action: make([]*metrics.Action, 101),
			},
			err: requestInvalidError,
		},
		// Case 003, nil action
		{
			req: &metrics.CounterI{
				Action: []*metrics.Action{
					{Metric: "valid", Number: 1},
					nil,
				},
			},
			err: actionInvalidError,
		},
		// Case 004, empty metric
		{
			req: &metrics.CounterI{
				Action: []*metrics.Action{{Metric: ""}},
			},
			err: actionInvalidError,
		},
		// Case 005, metric name too long
		{
			req: &metrics.CounterI{
				Action: []*metrics.Action{{Metric: string(make([]byte, 256))}},
			},
			err: actionInvalidError,
		},
		// Case 006, negative number
		{
			req: &metrics.CounterI{
				Action: []*metrics.Action{{Metric: "ok", Number: -1}},
			},
			err: actionInvalidError,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			err := verify(tc.req)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected %#v got %#v", tc.err, err)
			}
		})
	}
}

func Test_Server_Handler_Metrics_verify_success(t *testing.T) {
	testCases := []struct {
		req Request
		err error
	}{
		// Case 000
		{
			req: &metrics.CounterI{
				Action: []*metrics.Action{{Metric: "valid", Number: 1}},
			},
		},
		// Case 001
		{
			req: &metrics.CounterI{
				Action: []*metrics.Action{{Metric: "some_more_valid_counter", Number: 55.3, Labels: map[string]string{
					"foo": "bar",
					"baz": "zap",
				}}},
			},
		},
		// Case 002
		{
			req: &metrics.GaugeI{
				Action: []*metrics.Action{{Metric: "valid", Number: 2}},
			},
		},
		// Case 003
		{
			req: &metrics.GaugeI{
				Action: []*metrics.Action{{Metric: "some_more_valid_gauge", Number: 2635.0028, Labels: map[string]string{
					"hello": "world",
				}}},
			},
		},
		// Case 004
		{
			req: &metrics.HistogramI{
				Action: []*metrics.Action{{Metric: "valid", Number: 3}},
			},
		},
		// Case 005
		{
			req: &metrics.HistogramI{
				Action: []*metrics.Action{{Metric: "some_more_valid_histogram", Number: 11.0, Labels: map[string]string{
					"something": "something",
					"what":      "do",
					"you":       "know",
				}}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			err := verify(tc.req)
			if err != nil {
				t.Fatalf("expected %#v got %#v", nil, err)
			}
		})
	}
}
