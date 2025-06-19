package container

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/xh3b4sd/tracer"
)

type service struct {
	// hlt is the binary container health status, either 0 or 1.
	hlt float64
	// lab is the respective service label, e.g. alloy or specta.
	lab string
}

// service determines the binary container health of all ECS services as defined
// by the provided list of service details. The healthy status 1 is assigned to
// all services that have all of their desired containers running. Otherwise the
// unhealthy status 0 is assigned.
func (h *Handler) service(det []detail) ([]service, error) {
	var err error

	var ser []service
	for _, x := range det {
		var inp *ecs.DescribeServicesInput
		{
			inp = &ecs.DescribeServicesInput{
				Cluster:  aws.String(x.clu),
				Services: []string{x.arn},
				Include:  []types.ServiceField{types.ServiceFieldTags},
			}
		}

		var out *ecs.DescribeServicesOutput
		{
			out, err = h.ecs.DescribeServices(context.Background(), inp)
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		for _, y := range out.Services {
			var tag string
			{
				tag = serTag(y.Tags)
			}

			if tag == "" {
				h.log.Log(
					"level", "warning",
					"message", "skipping instrumentation for ECS service",
					"reason", "ECS service has no 'service' tag",
					"cluster", *y.ClusterArn,
					"service", *y.ServiceArn,
				)

				{
					continue
				}
			}

			var hlt float64 = 1
			if y.DesiredCount == 0 || y.RunningCount != y.DesiredCount {
				hlt = 0
			}

			ser = append(ser, service{
				hlt: hlt,
				lab: tag,
			})
		}
	}

	return ser, nil
}

func serTag(tag []types.Tag) string {
	for _, x := range tag {
		if *x.Key == "service" {
			return *x.Value
		}
	}

	return ""
}
