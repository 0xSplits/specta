package build

import (
	"context"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/tracer"
)

const (
	// max is the maximum amount of workflow executions to fetch at once. We never
	// have to look too far into the past, because we only ever have to observe
	// those builds that we have missed during either the most recent observation,
	// or the last deployment/downtime.
	max = 50
)

type workflow struct {
	// kin is the workflow type that we instrument here, e.g. check or image.
	kin string
	// lat is the duration that any given pipeline execution took to complete.
	lat time.Duration
	// rep is the repository label that we instrument here.
	rep string
	// rid is the run ID of the observed workflow.
	rid int64
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
func (h *Handler) workflow(det []detail) ([]workflow, error) {
	var err error

	var wor []workflow
	for _, x := range det {
		var inp *github.ListWorkflowRunsOptions
		{
			inp = &github.ListWorkflowRunsOptions{
				ListOptions: github.ListOptions{
					PerPage: max,
				},
			}
		}

		var out *github.WorkflowRuns
		{
			out, _, err = h.git.Actions.ListRepositoryWorkflowRuns(context.Background(), Organization, x.name, inp)
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		for _, y := range out.WorkflowRuns {
			wor = h.append(wor, x, y)
		}
	}

	return wor, nil
}

// append manages the skipping behaviour of already observed pipeline
// executions, as well as those pipeline executions that occured in the past
// before the internal start time.
func (h *Handler) append(wor []workflow, det detail, run *github.WorkflowRun) []workflow {
	// We skip all workflow executions that are too far in the past, given
	// Specta's own launch time.
	if run.GetCreatedAt().Before(h.sta) {
		return wor
	}

	// We skip all workflow executions that we already observed, given the cached
	// run ID.
	if h.cac.Contains(run.GetID()) {
		return wor
	}

	// We skip all workflow executions that have not yet completed, given the
	// workflow status string.
	if run.GetStatus() != "completed" {
		return wor
	}

	// We skip all workflow executions that are not whitelisted based on their
	// workflow name.
	if run.GetName() != det.check && run.GetName() != det.image {
		return wor
	}

	var lat time.Duration
	{

		lat = run.UpdatedAt.Sub(run.GetCreatedAt().Time)
	}

	var suc string
	if run.GetConclusion() == "success" {
		suc = "true"
	} else {
		suc = "false"
	}

	wor = append(wor, workflow{
		kin: worKin(det, run.GetName()),
		lat: lat,
		rep: det.label,
		rid: run.GetID(),
		suc: suc,
	})

	{
		h.cac.Add(run.GetID(), struct{}{})
	}

	return wor
}

func worKin(det detail, nam string) string {
	if nam == det.check {
		return "check"
	}

	if nam == det.image {
		return "image"
	}

	return ""
}
