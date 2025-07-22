package build

import (
	"context"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

const (
	// max is the maximum amount of workflow runs to fetch at once. We never have
	// to look too far into the past, because we only ever have to observe those
	// builds that we have missed during either the most recent observation, or
	// the last deployment/downtime. This number is a bit larger because
	// repositories like 0xSplits/splits contain many different workflows for that
	// monorepo, and in order to find the build runs that we are interested in, we
	// have to gather a broader picture for our workflows of interest not be be
	// "crowded out".
	max = 50
)

type workflow struct {
	// kin is the workflow type that we instrument here, e.g. check or image.
	kin string
	// lat is the duration that any given workflow execution took to complete.
	lat time.Duration
	// rep is the repository label that we instrument here.
	rep string
	// rid is the run ID of the observed workflow.
	rid int64
	// suc expresses whether the given deployment succeeded or failed, either true
	// or false.
	suc string
}

// workflow determines the latency of build runs for the provided workflow
// details. Build runs already observed are being skipped. Build runs with a
// start time before Specta launched are also being skipped. Both of those
// filters intent to prevent the duplication of latency measurements and their
// respective success/failure counts. With this stateless approach we prefer to
// rather miss a measurement than counting it twice, because counting twice
// creates inconsistencies in our SLOs.
func (h *Handler) workflow(det []detail) ([]workflow, error) {
	var wor []workflow
	var err error

	fnc := func(_ int, d detail) error {
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
			out, _, err = h.git.Actions.ListRepositoryWorkflowRuns(context.Background(), Organization, d.repo, inp)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		for _, x := range out.WorkflowRuns {
			wor = h.append(wor, d, x)
		}

		return nil
	}

	{
		err = parallel.Slice(det, fnc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return wor, nil
}

// append manages the skipping behaviour of already observed build runs, as well
// as those build runs that occured in the past before the internal start time.
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
