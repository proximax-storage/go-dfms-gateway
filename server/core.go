package server

import (
	"log"

	"github.com/valyala/fasthttp"
)

type gateway struct {
	server  fasthttp.Server
	address string
}

func NewGateway(opts ...GatewayOption) *gateway {
	gopts := ParseOptions(opts...)

	cfg, err := loadConfig(gopts.cfg)
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	gopts.ApplyToConfig(cfg)

	return &gateway{
		server: fasthttp.Server{
			Handler:      newMiddleware(newGatewayHandler(cfg.ApiAddress)),
			Name:         cfg.Name,
			GetOnly:      true,
			LogAllErrors: gopts.debug,
		},
		address: cfg.Address,
	}
}

func (g *gateway) Start() error {
	return g.server.ListenAndServe(g.address)
}

func (g *gateway) Stop() error {
	return g.server.Shutdown()
}
