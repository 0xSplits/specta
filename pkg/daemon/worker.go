package daemon

import (
	"github.com/0xSplits/specta/pkg/worker"
	"github.com/0xSplits/specta/pkg/worker/handler"
	"github.com/0xSplits/specta/pkg/worker/handler/privatekey"
)

func (d *Daemon) Worker() *worker.Worker {
	return worker.New(worker.Config{
		Han: []handler.Interface{
			privatekey.New(privatekey.Config{Env: d.env, Log: d.log}),
		},
		Log: d.log,
	})
}
