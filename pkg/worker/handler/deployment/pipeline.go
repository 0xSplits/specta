package deployment

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"github.com/xh3b4sd/tracer"
)

const (
	// max is the maximum amount of pipeline executions to fetch at once. We never
	// have to look too far into the past, because we only ever have to observe
	// those deployments that we have missed during either the most recent
	// observation, or the last deployment/downtime.
	max = 10
)

type pipeline struct {
	// eid is the execution ID of the observed pipeline.
	eid string
	// lat is the duration that any given pipeline execution took to complete.
	lat time.Duration
	// suc expresses whether the given deployment succeeded or failed, either true
	// or false.
	suc string
}

// pipeline determines the latency of CodePipeline executions for the provided
// pipeline details. Pipeline executions already observed are being skipped.
// Pipeline executions with a start time before Specta launched are also being
// skipped. Both of those filters intent to prevent the duplication of latency
// measurements and their respective success/failure counts. With this stateless
// approach we prefer to rather miss a measurement than counting it twice,
// because counting twice creates inconsistencies in our SLOs.
func (h *Handler) pipeline(det []detail) ([]pipeline, error) {
	var err error

	var pip []pipeline
	for _, x := range det {
		var inp *codepipeline.ListPipelineExecutionsInput
		{
			inp = &codepipeline.ListPipelineExecutionsInput{
				MaxResults:   aws.Int32(max),
				PipelineName: aws.String(x.nam),
			}
		}

		var out *codepipeline.ListPipelineExecutionsOutput
		{
			out, err = h.acp.ListPipelineExecutions(context.Background(), inp)
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		for _, y := range out.PipelineExecutionSummaries {
			pip = h.append(pip, y)
		}
	}

	return pip, nil
}

// append manages the skipping behaviour of already observed pipeline
// executions, as well as those pipeline executions that occured in the past
// before the internal start time.
func (h *Handler) append(pip []pipeline, sum types.PipelineExecutionSummary) []pipeline {
	// We skip all pipeline executions that are too far in the past, given
	// Specta's own launch time.
	if sum.StartTime.Before(h.sta) {
		return pip
	}

	// We skip all pipeline executions that we already observed, given the
	// cached execution ID.
	if h.cac.Contains(*sum.PipelineExecutionId) {
		return pip
	}

	// We skip all pipeline executions that have not yet completed, given the
	// pipeline status string.
	if pipSkp(sum.Status) {
		return pip
	}

	var lat time.Duration
	{
		lat = sum.LastUpdateTime.Sub(*sum.StartTime)
	}

	var suc string
	if sum.Status == types.PipelineExecutionStatusSucceeded {
		suc = "true"
	} else {
		suc = "false"
	}

	pip = append(pip, pipeline{
		eid: *sum.PipelineExecutionId,
		lat: lat,
		suc: suc,
	})

	{
		h.cac.Add(*sum.PipelineExecutionId, struct{}{})
	}

	return pip
}

func pipSkp(sta types.PipelineExecutionStatus) bool {
	return sta == types.PipelineExecutionStatusInProgress || sta == types.PipelineExecutionStatusStopping
}
