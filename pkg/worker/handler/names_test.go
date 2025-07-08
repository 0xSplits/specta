package handler

import (
	"fmt"
	"slices"
	"testing"

	"github.com/0xSplits/specta/pkg/worker/handler/build"
	"github.com/0xSplits/specta/pkg/worker/handler/deployment"
	"github.com/0xSplits/specta/pkg/worker/handler/endpoint"
	"github.com/0xSplits/specta/pkg/worker/handler/keypair"
	"github.com/0xSplits/specta/pkg/worker/handler/stack"
)

func Test_Worker_Handler_Names(t *testing.T) {
	testCases := []struct {
		han []Interface
		nam []string
	}{
		// Case 000
		{
			han: []Interface{
				&build.Handler{},
				&deployment.Handler{},
				&endpoint.Handler{},
				&keypair.Handler{},
				&stack.Handler{},
			},
			nam: []string{
				"build",
				"deployment",
				"endpoint",
				"keypair",
				"stack",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			nam := Names(tc.han)
			if !slices.Equal(nam, tc.nam) {
				t.Fatalf("expected %#v got %#v", tc.nam, nam)
			}
		})
	}
}
