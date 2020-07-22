package server

import (
	apihttp "github.com/proximax-storage/go-xpx-dfms-api-http"

	"github.com/valyala/fasthttp"
)

type gatewayServer struct {
	server  fasthttp.Server
	address string
}

func NewGatewayServer(cfg *GatewayConfig) *gatewayServer {
	return &gatewayServer{
		server: fasthttp.Server{
			Handler:      NewGatewayHandler(apihttp.NewClientAPI(cfg.AddressAPI)).Serve,
			Name:         cfg.Name,
			GetOnly:      cfg.GetOnly,
			LogAllErrors: cfg.LogError,
		},
		address: cfg.Address,
	}
}

func DefaultGatewayServer() *gatewayServer {
	return NewGatewayServer(DefaultGatewayConfig())
}

func (g *gatewayServer) Start() error {
	return g.server.ListenAndServe(g.address)
}

func (g *gatewayServer) Stop() error {
	return g.server.Shutdown()
}