package stack

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

var status = map[types.ResourceStatus]float64{
	types.ResourceStatusCreateFailed:             0.0,
	types.ResourceStatusDeleteFailed:             0.0,
	types.ResourceStatusDeleteSkipped:            0.0,
	types.ResourceStatusExportFailed:             0.0,
	types.ResourceStatusExportRollbackComplete:   0.0,
	types.ResourceStatusExportRollbackFailed:     0.0,
	types.ResourceStatusExportRollbackInProgress: 0.0,
	types.ResourceStatusImportFailed:             0.0,
	types.ResourceStatusImportRollbackComplete:   0.0,
	types.ResourceStatusImportRollbackFailed:     0.0,
	types.ResourceStatusImportRollbackInProgress: 0.0,
	types.ResourceStatusRollbackComplete:         0.0,
	types.ResourceStatusRollbackFailed:           0.0,
	types.ResourceStatusRollbackInProgress:       0.0,
	types.ResourceStatusUpdateFailed:             0.0,
	types.ResourceStatusUpdateRollbackComplete:   0.0,
	types.ResourceStatusUpdateRollbackFailed:     0.0,
	types.ResourceStatusUpdateRollbackInProgress: 0.0,

	types.ResourceStatusCreateInProgress: 0.5,
	types.ResourceStatusDeleteInProgress: 0.5,
	types.ResourceStatusExportInProgress: 0.5,
	types.ResourceStatusImportInProgress: 0.5,
	types.ResourceStatusUpdateInProgress: 0.5,

	types.ResourceStatusCreateComplete: 1.0,
	types.ResourceStatusDeleteComplete: 1.0,
	types.ResourceStatusExportComplete: 1.0,
	types.ResourceStatusImportComplete: 1.0,
	types.ResourceStatusUpdateComplete: 1.0,
}

type stack struct {
	// hlt is the CloudFormation stack status, either 0.0, 0.5 or 1.0.
	hlt float64
	// lab is the respective stack label, e.g. root or cache.
	lab string
}

// stack determines the CloudFormation stack health of all CloudFormation stacks
// as defined by the provided list of stack details. The healthy status 1 is
// assigned to all stacks that have a stack status suffix of _COMPLETE. The
// exception here are _ROLLBACK_COMPLETE statuses. The _PROGRESS suffix is
// assigned the stack status 0.5, otherwise the unhealthy stack status 0 is
// assigned.
func (h *Handler) stack(det []detail) ([]stack, error) {
	var err error

	var sta []stack
	for _, x := range det {
		var inp *cloudformation.DescribeStackResourcesInput
		{
			inp = &cloudformation.DescribeStackResourcesInput{
				StackName: aws.String(x.arn),
			}
		}

		var out *cloudformation.DescribeStackResourcesOutput
		{
			out, err = h.cfc.DescribeStackResources(context.Background(), inp)
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		for _, y := range out.StackResources {
			var tag string
			{
				tag = staTag(*y.LogicalResourceId)
			}

			if tag == "" {
				h.log.Log(
					"level", "warning",
					"message", "skipping instrumentation for CloudFormation stack",
					"reason", "CloudFormation stack name pattern is unrecognizable",
					"name", *y.LogicalResourceId,
				)

				{
					continue
				}
			}

			var hlt float64
			{
				hlt = status[y.ResourceStatus]
			}

			sta = append(sta, stack{
				hlt: hlt,
				lab: tag,
			})
		}
	}

	return sta, nil
}

func staTag(nam string) string {
	for k, v := range mapping {
		if strings.Contains(nam, k) {
			return v
		}
	}

	return ""
}
