package go_dfms_gateway

import (
	logging "github.com/ipfs/go-log"
	api "github.com/proximax-storage/go-xpx-dfms-api"

	"github.com/valyala/fasthttp"
)

var log = logging.Logger("gateway")

func init() {
	logging.SetupLogging()
}

type gateway struct {
	server  fasthttp.Server
	address string
}

func NewGateway(api api.Client, opts ...GatewayOption) *gateway {
	gopts := ParseOptions(opts...)

	cfg, err := loadConfig(gopts.cfg)
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	gopts.ApplyToConfig(cfg)

	if gopts.debug {
		err := logging.SetLogLevel("DEBUG", "gateway")
		if err != nil {
			log.Warn("Cannot load config: ", err)
		}
	}

	return &gateway{
		server: fasthttp.Server{
			Handler:      newMiddleware(newGatewayHandler(api)),
			Name:         cfg.Name,
			GetOnly:      cfg.GetOnly,
			LogAllErrors: gopts.debug,
		},
		address: cfg.Address,
	}
}

func (g *gateway) Start() error {
	println("Gateway listening at", g.address)
	return g.server.ListenAndServe(g.address)
}

func (g *gateway) Stop() error {
	println("Stopping gateway...")
	return g.server.Shutdown()
}
