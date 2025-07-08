package deployment

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/specta/pkg/runtime"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"github.com/xh3b4sd/logger"
)

func Test_Worker_Handler_Deployment_append(t *testing.T) {
	testCases := []struct {
		sta time.Time
		pip []pipeline
		sum types.PipelineExecutionSummary
		app []pipeline
	}{
		// Case 000, empty, append
		{
			sta: time.Unix(3, 0),
			pip: []pipeline{},
			sum: types.PipelineExecutionSummary{
				PipelineExecutionId: aws.String("foo"),
				StartTime:           aws.Time(time.Unix(5, 0)),
				LastUpdateTime:      aws.Time(time.Unix(8, 0)),
				Status:              types.PipelineExecutionStatusSucceeded,
			},
			app: []pipeline{
				{eid: "foo", lat: 3 * time.Second, suc: "true"},
			},
		},
		// Case 001, "foo" already cached, skip
		{
			sta: time.Unix(3, 0),
			pip: []pipeline{
				{eid: "foo", lat: 3 * time.Second, suc: "true"},
			},
			sum: types.PipelineExecutionSummary{
				PipelineExecutionId: aws.String("foo"),
				StartTime:           aws.Time(time.Unix(5, 0)),
				LastUpdateTime:      aws.Time(time.Unix(8, 0)),
				Status:              types.PipelineExecutionStatusSucceeded,
			},
			app: []pipeline{
				{eid: "foo", lat: 3 * time.Second, suc: "true"},
			},
		},
		// Case 002, empty, before start time, skip
		{
			sta: time.Unix(6, 0),
			pip: []pipeline{},
			sum: types.PipelineExecutionSummary{
				PipelineExecutionId: aws.String("foo"),
				StartTime:           aws.Time(time.Unix(5, 0)),
				LastUpdateTime:      aws.Time(time.Unix(8, 0)),
				Status:              types.PipelineExecutionStatusSucceeded,
			},
			app: []pipeline{},
		},
		// Case 003, still running, skip
		{
			sta: time.Unix(3, 0),
			pip: []pipeline{},
			sum: types.PipelineExecutionSummary{
				PipelineExecutionId: aws.String("foo"),
				StartTime:           aws.Time(time.Unix(5, 0)),
				LastUpdateTime:      aws.Time(time.Unix(8, 0)),
				Status:              types.PipelineExecutionStatusInProgress,
			},
			app: []pipeline{},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var han *Handler
			{
				han = New(Config{
					Aws: aws.Config{Region: "us-west-2"},
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
			// handler. So we add the input execution IDs to the real cache
			// implementation in order to verify the skipping process.
			for _, x := range tc.pip {
				han.cac.Add(x.eid, struct{}{})
			}

			app := han.append(tc.pip, tc.sum)
			if !slices.Equal(app, tc.app) {
				t.Fatalf("expected %#v got %#v", tc.app, app)
			}
		})
	}
}
