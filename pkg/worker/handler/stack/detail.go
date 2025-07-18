package stack

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

type detail struct {
	// nam is the name of the CloudFormation stack.
	nam string
	// sta is the raw CloudFormation stack status, e.g. UPDATE_COMPLETE.
	sta types.StackStatus
}

// detail finds all CloudFormation stack statuses that are tagged with the
// "environment" that matches Specta's runtime configuration. In other words, if
// Specta is running in "staging", then detail() will find all CloudFormation
// stacks labelled with the resource tags environment=staging.
func (h *Handler) detail() ([]detail, error) {
	var err error

	// Create a new paginator for each instrumentation cycle in order to guarantee
	// data integrity per loop.
	var pag *cloudformation.DescribeStacksPaginator
	{
		pag = cloudformation.NewDescribeStacksPaginator(h.cfc, &cloudformation.DescribeStacksInput{})
	}

	var det []detail
	for pag.HasMorePages() {
		var out *cloudformation.DescribeStacksOutput
		{
			out, err = pag.NextPage(context.Background())
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		for _, x := range out.Stacks {
			if !hasEnv(x.Tags, h.env.Environment) {
				continue
			}

			det = append(det, detail{
				nam: *x.StackName,
				sta: x.StackStatus,
			})
		}
	}

	if len(det) == 0 {
		return nil, tracer.Mask(missingRootStackError, tracer.Context{Key: "environment", Value: h.env.Environment})
	}

	return det, nil
}

func hasEnv(tags []types.Tag, env string) bool {
	for _, t := range tags {
		if t.Key != nil && t.Value != nil && *t.Key == "environment" && *t.Value == env {
			return true
		}
	}

	return false
}
