package deployment

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	tagtypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/xh3b4sd/tracer"
)

type detail struct {
	// nam is the name of the AWS CodePipeline to instrument for the current
	// environment.
	nam string
}

// detail finds all AWS CodePipeline names that are tagged with the
// "environment" that matches Specta's runtime configuration. In other words, if
// Specta is running in "staging", then detail() will find all AWS CodePipeline
// names labelled with the resource tags environment=staging.
func (h *Handler) detail() ([]detail, error) {
	var err error

	var inp *resourcegroupstaggingapi.GetResourcesInput
	{
		inp = &resourcegroupstaggingapi.GetResourcesInput{
			ResourceTypeFilters: []string{"codepipeline:pipeline"},
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
		var spl []string
		{
			spl = strings.Split(*x.ResourceARN, ":")
		}

		var nam string
		{
			nam = spl[len(spl)-1] // the last element divided by a colon is the pipeline name
		}

		det = append(det, detail{
			nam: nam,
		})
	}

	return det, nil
}
