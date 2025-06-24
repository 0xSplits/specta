package stack

import (
	"fmt"
	"testing"
)

func Test_Worker_Handler_Stack_rooSta(t *testing.T) {
	testCases := []struct {
		arn string
		roo bool
	}{
		// Case 000
		{
			arn: "server-test-SpectaStack-XQCAJRTZL8RF",
			roo: false,
		},
		// Case 001
		{
			arn: "server-test-DiscoveryStack-1HJ4KRXNU9IZZ",
			roo: false,
		},
		// Case 002
		{
			arn: "server-test-CacheStack-1CLF4ZWT6S9EZ",
			roo: false,
		},
		// Case 003
		{
			arn: "server-test-FargateStack-QGXQ9XZ4J44K",
			roo: false,
		},
		// Case 004
		{
			arn: "server-test-VpcStack-P4032W206SOD",
			roo: false,
		},
		// Case 005
		{
			arn: "server-test-TelemetryStack-1CI1X2G5NOG2J",
			roo: false,
		},
		// Case 006
		{
			arn: "server-test-RdsStack-TF2N40EHYUUN",
			roo: false,
		},
		// Case 007
		{
			arn: "server-test",
			roo: true,
		},
		// Case 008
		{
			arn: "another-root-stack-name",
			roo: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			roo := rooSta(tc.arn)
			if roo != tc.roo {
				t.Fatalf("expected %#v got %#v", tc.roo, roo)
			}
		})
	}
}
