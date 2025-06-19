package container

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	tagtypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/xh3b4sd/tracer"
)

type detail struct {
	// arn is the well defined Amazon Resource Name of the ECS service.
	arn string
	// clu is the short cluster name that the given service is part of.
	clu string
}

// detail finds all ECS service ARNs that are tagged with the "environment" that
// matches Specta's runtime configuration. In other words, if Specta is running
// in "staging", then detail() will find all ECS services labelled with the
// resource tags environment=staging.
func (h *Handler) detail() ([]detail, error) {
	var err error

	var inp *resourcegroupstaggingapi.GetResourcesInput
	{
		inp = &resourcegroupstaggingapi.GetResourcesInput{
			ResourceTypeFilters: []string{"ecs:service"},
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
		spl := strings.Split(*x.ResourceARN, "/")
		if len(spl) != 3 {
			return nil, tracer.Maskf(invalidAmazonResourceNameError, "%s", *x.ResourceARN)
		}

		det = append(det, detail{
			arn: *x.ResourceARN,
			clu: spl[len(spl)-2],
		})
	}

	return det, nil
}
