package stack

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

var (
	// exp defines a regular expression for the stack names of CloudFormation
	// stacks. Note that there is no precise definition for the amount of
	// characters that a random hash suffix for the stack names ought to look
	// like. So we are testing against at least 10 characters, but also accept
	// e.g. 13 or 14.
	//
	//     server-test-FargateStack-QGXQ9XZ4J44K
	//
	exp = regexp.MustCompile(`-[A-Z0-9]{10,}$`)

	status = map[types.StackStatus]float64{
		types.StackStatusCreateFailed:                            0.0,
		types.StackStatusDeleteFailed:                            0.0,
		types.StackStatusImportRollbackComplete:                  0.0,
		types.StackStatusImportRollbackFailed:                    0.0,
		types.StackStatusImportRollbackInProgress:                0.0,
		types.StackStatusRollbackComplete:                        0.0,
		types.StackStatusRollbackFailed:                          0.0,
		types.StackStatusRollbackInProgress:                      0.0,
		types.StackStatusUpdateFailed:                            0.0,
		types.StackStatusUpdateRollbackComplete:                  0.0,
		types.StackStatusUpdateRollbackCompleteCleanupInProgress: 0.0,
		types.StackStatusUpdateRollbackFailed:                    0.0,
		types.StackStatusUpdateRollbackInProgress:                0.0,

		types.StackStatusCreateInProgress:                0.5,
		types.StackStatusDeleteInProgress:                0.5,
		types.StackStatusImportInProgress:                0.5,
		types.StackStatusReviewInProgress:                0.5,
		types.StackStatusUpdateCompleteCleanupInProgress: 0.5,
		types.StackStatusUpdateInProgress:                0.5,

		types.StackStatusCreateComplete: 1.0,
		types.StackStatusDeleteComplete: 1.0,
		types.StackStatusImportComplete: 1.0,
		types.StackStatusUpdateComplete: 1.0,
	}
)

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
	var sta []stack

	for _, x := range det {
		var tag string
		{
			tag = staTag(x.nam)
		}

		if tag == "" {
			h.log.Log(
				"level", "warning",
				"message", "skipping instrumentation for CloudFormation stack",
				"reason", "CloudFormation stack name pattern is unrecognizable",
				"name", x.nam,
			)

			{
				continue
			}
		}

		var hlt float64
		{
			hlt = status[x.sta]
		}

		sta = append(sta, stack{
			hlt: hlt,
			lab: tag,
		})
	}

	return sta, nil
}

// rooSta identifies whether the given CloudFormation stack name is considered a
// root stack name.
func rooSta(arn string) bool {
	return !exp.MatchString(arn)
}

func staTag(nam string) string {
	for k, v := range mapping {
		if strings.Contains(nam, k) {
			return v
		}
	}

	if rooSta(nam) {
		return "root"
	}

	return ""
}
