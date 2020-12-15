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
	server     fasthttp.Server
	address    string
	enableCORs bool
	cors       *cors
}

func NewGateway(api api.Client, opts ...GatewayOption) *gateway {
	gopts := ParseOptions(opts...)

	cfg, err := loadConfig(gopts.cfg)
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	gopts.ApplyToConfig(cfg)

	if gopts.debug {
		err := logging.SetLogLevel("gateway", "DEBUG")
		if err != nil {
			println("Cannot set DEBUG mode for gateway: ", err.Error())
		}
	}

	handler := newMiddleware(newGatewayHandler(api), cfg.CORs)
	return &gateway{
		server: fasthttp.Server{
			Handler:      handler,
			Name:         cfg.Name,
			GetOnly:      cfg.GetOnly,
			LogAllErrors: gopts.debug,
		},
		address:    cfg.Address,
		enableCORs: cfg.CORs.Enable,
		cors:       cfg.CORs,
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
