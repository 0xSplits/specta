package daemon

import (
	"github.com/0xSplits/specta/pkg/server"
	"github.com/0xSplits/workit/worker/parallel"
)

type Interface interface {
	Server() *server.Server
	Worker() *parallel.Worker
}
