package server

import (
	"os"
	"os/signal"
	"syscall"

	logging "github.com/ipfs/go-log"

	"github.com/valyala/fasthttp"
)

var log = logging.Logger("core")

func init() {
	logging.SetupLogging()
}

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

	if gopts.debug {
		logging.SetDebugLogging()
	}

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
	println("Gateway listening at", g.address)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs

		err := g.Stop()
		if err != nil {
			log.Error(err)
		}
	}()

	return g.server.ListenAndServe(g.address)
}

func (g *gateway) Stop() error {
	println("Stopping gateway...")
	return g.server.Shutdown()
}
