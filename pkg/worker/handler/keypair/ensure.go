package keypair

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	var out *ec2.DescribeKeyPairsOutput
	{
		out, err = h.ec2.DescribeKeyPairs(context.Background(), nil)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err = h.reg.Counter("aws_ec2_keypair_total", float64(len(out.KeyPairs)), nil)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
