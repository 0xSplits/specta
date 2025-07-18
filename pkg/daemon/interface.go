package daemon

import (
	"github.com/0xSplits/specta/pkg/server"
	"github.com/0xSplits/workit/engine"
)

type Interface interface {
	Server() *server.Server
	Worker() *engine.Engine
}
