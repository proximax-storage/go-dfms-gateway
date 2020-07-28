package server

import (
	apihttp "github.com/proximax-storage/go-xpx-dfms-api-http"

	"github.com/valyala/fasthttp"
)

type gateway struct {
	server  fasthttp.Server
	address string
}

func NewGateway(cfg *GatewayConfig) *gateway {
	handler := newMiddleware(newGatewayHandler(apihttp.NewClientAPI(cfg.AddressAPI)))
	return &gateway{
		server: fasthttp.Server{
			Handler:      handler,
			Name:         cfg.Name,
			GetOnly:      cfg.GetOnly,
			LogAllErrors: cfg.LogError,
		},
		address: cfg.Address,
	}
}

func DefaultGateway() *gateway {
	return NewGateway(DefaultGatewayConfig())
}

func (g *gateway) Start() error {
	return g.server.ListenAndServe(g.address)
}

func (g *gateway) Stop() error {
	return g.server.Shutdown()
}
