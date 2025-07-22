package container

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

type service struct {
	// hlt is the container health status, either 0.0, 0.5 or 1.0.
	hlt float64
	// lab is the respective service label, e.g. alloy or specta.
	lab string
}

// service determines the binary container health of all ECS services as defined
// by the provided list of service details. The healthy status 1 is assigned to
// all services that have all of their desired containers running. Otherwise the
// unhealthy status 0 is assigned.
func (h *Handler) service(det []detail) ([]service, error) {
	var ser []service
	var err error

	fnc := func(_ int, d detail) error {
		var inp *ecs.DescribeServicesInput
		{
			inp = &ecs.DescribeServicesInput{
				Cluster:  aws.String(d.clu),
				Services: []string{d.arn},
				Include:  []types.ServiceField{types.ServiceFieldTags},
			}
		}

		var out *ecs.DescribeServicesOutput
		{
			out, err = h.ecs.DescribeServices(context.Background(), inp)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		for _, x := range out.Services {
			ser = h.append(ser, x)
		}

		return nil
	}

	{
		err = parallel.Slice(det, fnc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return ser, nil
}

func (h *Handler) append(ser []service, out types.Service) []service {
	var tag string
	{
		tag = serTag(out.Tags)
	}

	if tag == "" {
		h.log.Log(
			"level", "warning",
			"message", "skipping instrumentation for ECS service",
			"reason", "ECS service has no 'service' tag",
			"cluster", *out.ClusterArn,
			"service", *out.ServiceArn,
		)

		{
			return ser
		}
	}

	var hlt float64
	switch {
	case out.RunningCount == 0:
		hlt = 0 // no containers running
	case out.RunningCount != out.DesiredCount:
		hlt = 0.5 // not enough containers running
	default:
		hlt = 1 // all containers running
	}

	ser = append(ser, service{
		hlt: hlt,
		lab: tag,
	})

	return ser
}

func serTag(tag []types.Tag) string {
	for _, x := range tag {
		if *x.Key == "service" {
			return *x.Value
		}
	}

	return ""
}
