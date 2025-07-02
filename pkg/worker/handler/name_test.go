package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/0xSplits/specta/pkg/worker/handler/container"
	"github.com/0xSplits/specta/pkg/worker/handler/keypair"
)

func Test_Worker_Handler_Name(t *testing.T) {
	testCases := []struct {
		han Interface
		nam string
	}{
		// Case 000
		{
			han: &testHandler{},
			nam: "handler",
		},
		// Case 001
		{
			han: (*testHandler)(nil),
			nam: "handler",
		},
		// Case 002
		{
			han: &container.Handler{},
			nam: "container",
		},
		// Case 003
		{
			han: &keypair.Handler{},
			nam: "keypair",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			nam := Name(tc.han)
			if nam != tc.nam {
				t.Fatalf("expected %#v got %#v", tc.nam, nam)
			}
		})
	}
}

type testHandler struct{}

func (h *testHandler) Cooler() time.Duration {
	return 0
}

func (h *testHandler) Ensure() error {
	return nil
}
