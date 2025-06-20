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
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-SpectaStack-XQCAJRTZL8RF/9b0d7690-491f-11f0-9d9a-02e5e5b98f6f",
			roo: false,
		},
		// Case 001
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-DiscoveryStack-1HJ4KRXNU9IZZ/7643e420-491f-11f0-8f62-0631b64352f1",
			roo: false,
		},
		// Case 002
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-CacheStack-1CLF4ZWT6S9EZ/0b6f86e0-30ce-11f0-ad82-0a7b61ec7c77",
			roo: false,
		},
		// Case 003
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-FargateStack-QGXQ9XZ4J44K/165deb30-30d0-11f0-9eeb-023c6b26cb57",
			roo: false,
		},
		// Case 004
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-VpcStack-P4032W206SOD/9e159990-30cd-11f0-9975-061fb88c06d9",
			roo: false,
		},
		// Case 005
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-TelemetryStack-1CI1X2G5NOG2J/0399cdc0-34bd-11f0-bd75-061643229203",
			roo: false,
		},
		// Case 006
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test-RdsStack-TF2N40EHYUUN/0b61a430-30ce-11f0-9e0a-06c6e436b1f3",
			roo: false,
		},
		// Case 007
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/server-test/9c95fe70-30cd-11f0-b8b6-0a489308d945",
			roo: true,
		},
		// Case 008
		{
			arn: "arn:aws:cloudformation:us-west-2:995626699990:stack/another-root-stack-name/9c95fe70-30cd-11f0-b8b6-0a489308d945",
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
