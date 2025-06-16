package registry

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/xh3b4sd/logger"
)

func Test_Registry_record_counter_failure(t *testing.T) {
	testCases := []struct {
		nam string
		val float64
		lab map[string]string
		err error
	}{
		// Case 000, metric name not registered
		{
			nam: "not_allowed",
			val: 38.5,
			lab: map[string]string{
				"foo": "one",
			},
			err: metricNameWhitelistError,
		},
		// Case 001, label key not registered
		{
			nam: "allowed_counter",
			val: 38.5,
			lab: map[string]string{
				"not-allowed": "two",
			},
			err: labelKeyWhitelistError,
		},
		// Case 002, label value not registered
		{
			nam: "allowed_counter",
			val: 38.5,
			lab: map[string]string{
				"foo": "bar",
			},
			err: labelValueWhitelistError,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var rec *recorder.Fake
			{
				rec = recorder.NewFake(recorder.FakeConfig{
					Lab: map[string][]string{
						"foo": {"one", "two"},
					},
				})
			}

			var reg Interface
			{
				reg = New(Config{
					Env: envvar.Env{
						Environment: "testing",
					},
					Log: logger.Fake(),

					Cou: map[string]recorder.Interface{
						"allowed_counter": rec,
					},
					Gau: map[string]recorder.Interface{},
					His: map[string]recorder.Interface{},
				})
			}

			err := reg.Counter(tc.nam, tc.val, tc.lab)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected %#v got %#v", tc.err, err)
			}

			if rec.Recorded().Lab != nil {
				t.Fatalf("expected %#v got %#v", nil, rec.Recorded().Lab)
			}
			if rec.Recorded().Val != nil {
				t.Fatalf("expected %#v got %#v", nil, rec.Recorded().Val)
			}
		})
	}
}

func Test_Registry_record_counter_success(t *testing.T) {
	testCases := []struct {
		nam string
		val float64
		lab map[string]string
	}{
		// Case 000
		{
			nam: "allowed_counter",
			val: 38.5,
			lab: map[string]string{
				"foo": "one",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var rec *recorder.Fake
			{
				rec = recorder.NewFake(recorder.FakeConfig{
					Lab: map[string][]string{
						"foo": {"one", "two"},
					},
				})
			}

			var reg Interface
			{
				reg = New(Config{
					Env: envvar.Env{
						Environment: "testing",
					},
					Log: logger.Fake(),

					Cou: map[string]recorder.Interface{
						"allowed_counter": rec,
					},
					Gau: map[string]recorder.Interface{},
					His: map[string]recorder.Interface{},
				})
			}

			err := reg.Counter(tc.nam, tc.val, tc.lab)
			if err != nil {
				t.Fatalf("expected %#v got %#v", nil, err)
			}

			if len(rec.Recorded().Lab) != 1 {
				t.Fatalf("expected %#v got %#v", 1, len(rec.Recorded().Lab))
			}
			if !reflect.DeepEqual(rec.Recorded().Lab[0], map[string]string{"env": "testing", "foo": "one"}) {
				t.Fatalf("expected %#v got %#v", map[string]string{"env": "testing", "foo": "one"}, rec.Recorded().Lab[0])
			}

			if len(rec.Recorded().Val) != 1 {
				t.Fatalf("expected %#v got %#v", 1, len(rec.Recorded().Val))
			}
			if !reflect.DeepEqual(rec.Recorded().Val[0], 38.5) {
				t.Fatalf("expected %#v got %#v", 38.5, rec.Recorded().Val[0])
			}
		})
	}
}
