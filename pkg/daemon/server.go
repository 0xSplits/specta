package daemon

import (
	"net"

	"github.com/0xSplits/specta/pkg/server"
	"github.com/0xSplits/specta/pkg/server/handler"
	"github.com/0xSplits/specta/pkg/server/handler/metrics"
	"github.com/0xSplits/specta/pkg/server/interceptor/failure"
	"github.com/0xSplits/specta/pkg/server/middleware/cors"
	"github.com/gorilla/mux"
	"github.com/twitchtv/twirp"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Server() *server.Server {
	var err error

	var lis net.Listener
	{
		lis, err = net.Listen("tcp", net.JoinHostPort(d.env.HttpHost, d.env.HttpPort))
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return server.New(server.Config{
		Han: []handler.Interface{
			metrics.New(metrics.Config{Env: d.env, Log: d.log, Reg: d.reg}),
		},
		Int: []twirp.Interceptor{
			failure.New(failure.Config{Log: d.log}).Method,
		},
		Lis: lis,
		Log: d.log,
		Mid: []mux.MiddlewareFunc{
			cors.New(cors.Config{Log: d.log}).Handler,
		},
	})
}
