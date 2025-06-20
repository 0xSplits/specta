package stack

import (
	"context"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	tagtypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/xh3b4sd/tracer"
)

var (
	// exp defines a regular expression for the ARN suffixes of nested
	// CloudFormation stacks. Note that there is no precise definition for the
	// amount of characters that a random hack suffix for the stack names ought to
	// look like. So we are testing against at least 10 characters, but also
	// accept e.g. 13 or 14.
	//
	//     arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-FargateStack-QGXQ9XZ4J44K/165deb30-30d0-11f0-9eeb-023c6b26cb57
	//
	exp = regexp.MustCompile(`-[A-Z0-9]{10,}/[a-f0-9-]{36}$`)
)

type detail struct {
	// arn is the well defined Amazon Resource Name of the CloudFormation stack.
	arn string
}

// detail finds all CloudFormation stack ARNs that are tagged with the
// "environment" that matches Specta's runtime configuration. In other words, if
// Specta is running in "staging", then detail() will find all CloudFormation
// stacks labelled with the resource tags environment=staging.
func (h *Handler) detail() ([]detail, error) {
	var err error

	var inp *resourcegroupstaggingapi.GetResourcesInput
	{
		inp = &resourcegroupstaggingapi.GetResourcesInput{
			ResourceTypeFilters: []string{"cloudformation:stack"},
			TagFilters: []tagtypes.TagFilter{
				{
					Key:    aws.String("environment"),
					Values: []string{h.env.Environment},
				},
			},
		}
	}

	var out *resourcegroupstaggingapi.GetResourcesOutput
	{
		out, err = h.tag.GetResources(context.Background(), inp)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var det []detail
	for _, x := range out.ResourceTagMappingList {
		var roo bool
		{
			roo = rooSta(*x.ResourceARN)
		}

		if !roo {
			continue
		}

		det = append(det, detail{
			arn: *x.ResourceARN,
		})
	}

	if len(det) == 0 {
		return nil, tracer.Maskf(missingRootStackError, "%v", allArn(out.ResourceTagMappingList))
	}
	if len(det) > 1 {
		return nil, tracer.Maskf(tooManyRootStacksError, "%v", det)
	}

	return det, nil
}

func allArn(lis []tagtypes.ResourceTagMapping) []string {
	var all []string

	for _, x := range lis {
		all = append(all, *x.ResourceARN)
	}

	return all
}

// rooSta identifies whether the given CloudFormation stack ARN is considered a
// root stack ARN.
func rooSta(arn string) bool {
	return !exp.MatchString(arn)
}
