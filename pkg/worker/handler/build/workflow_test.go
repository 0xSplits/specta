package build

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/specta/pkg/runtime"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/logger"
)

func Test_Worker_Handler_Build_append(t *testing.T) {
	testCases := []struct {
		sta time.Time
		wor []workflow
		det detail
		run *github.WorkflowRun
		app []workflow
	}{
		// Case 000, empty, append
		{
			sta: time.Unix(3, 0),
			wor: []workflow{},
			det: detail{
				label: "server",
				check: "foo",
				image: "bar",
			},
			run: &github.WorkflowRun{
				ID:         github.Ptr(int64(47)),
				Name:       github.Ptr("foo"),
				CreatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(5, 0)}),
				UpdatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(8, 0)}),
				Status:     github.Ptr("completed"),
				Conclusion: github.Ptr("success"),
			},
			app: []workflow{{kin: "check", lat: 3 * time.Second, rep: "server", rid: 47, suc: "true"}},
		},
		// Case 001, 47 already cached, skip
		{
			sta: time.Unix(3, 0),
			wor: []workflow{{lat: 3 * time.Second, rep: "server", rid: 47, suc: "true"}},
			det: detail{
				label: "server",
				check: "foo",
				image: "bar",
			},
			run: &github.WorkflowRun{
				ID:         github.Ptr(int64(47)),
				Name:       github.Ptr("bar"),
				CreatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(5, 0)}),
				UpdatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(8, 0)}),
				Status:     github.Ptr("completed"),
				Conclusion: github.Ptr("success"),
			},
			app: []workflow{{lat: 3 * time.Second, rep: "server", rid: 47, suc: "true"}},
		},
		// Case 002, empty, before start time, skip
		{
			sta: time.Unix(6, 0),
			wor: nil,
			det: detail{
				label: "server",
				check: "foo",
				image: "bar",
			},
			run: &github.WorkflowRun{
				ID:         github.Ptr(int64(47)),
				Name:       github.Ptr("foo"),
				CreatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(5, 0)}),
				UpdatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(8, 0)}),
				Status:     github.Ptr("completed"),
				Conclusion: github.Ptr("success"),
			},
			app: nil,
		},
		// Case 003, still running, skip
		{
			sta: time.Unix(3, 0),
			wor: nil,
			det: detail{
				label: "server",
				check: "foo",
				image: "bar",
			},
			run: &github.WorkflowRun{
				ID:        github.Ptr(int64(47)),
				Name:      github.Ptr("bar"),
				CreatedAt: github.Ptr(github.Timestamp{Time: time.Unix(5, 0)}),
				UpdatedAt: github.Ptr(github.Timestamp{Time: time.Unix(8, 0)}),
				Status:    github.Ptr("in_progress"),
			},
			app: nil,
		},
		// Case 004, unsupported workflow, skip
		{
			sta: time.Unix(3, 0),
			wor: []workflow{},
			det: detail{
				label: "server",
				check: "foo",
				image: "bar",
			},
			run: &github.WorkflowRun{
				ID:         github.Ptr(int64(47)),
				Name:       github.Ptr("unsupported"),
				CreatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(5, 0)}),
				UpdatedAt:  github.Ptr(github.Timestamp{Time: time.Unix(8, 0)}),
				Status:     github.Ptr("completed"),
				Conclusion: github.Ptr("success"),
			},
			app: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var han *Handler
			{
				han = New(Config{
					Env: envvar.Env{Environment: "testing"},
					Log: logger.Fake(),
					Met: recorder.NewMeter(recorder.MeterConfig{
						Env: "testing",
						Ver: runtime.Tag(),
					}),
				})
			}

			// We have to control the start time that the handler is using internally,
			// so that we can verify teh time based skipping behaviour.
			{
				han.sta = tc.sta
			}

			// We have to simulate that the test input was already observed by the
			// handler. So we add the input run IDs to the real cache implementation
			// in order to verify the skipping process.
			for _, x := range tc.wor {
				han.cac.Add(x.rid, struct{}{})
			}

			app := han.append(tc.wor, tc.det, tc.run)
			if !slices.Equal(app, tc.app) {
				t.Fatalf("expected %#v got %#v", tc.app, app)
			}
		})
	}
}
