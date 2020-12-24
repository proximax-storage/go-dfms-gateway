package go_dfms_gateway

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func NewCors(allowedMethods []string, allowedHeaders []string, allowedOrigins []string) *cors {
	return &cors{
		Enable:         true,
		AllowedMethods: allowedMethods,
		AllowedHeaders: allowedHeaders,
		AllowedOrigins: allowedOrigins,
	}
}

func DefaultCors() *cors {
	return &cors{
		Enable:         false,
		AllowedMethods: []string{"GET", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"Accept", "Accept-Language", "Content-Language", "Content-Type"},
		AllowedOrigins: []string{"*"},
	}
}

type cors struct {
	Enable         bool
	AllowedMethods []string
	AllowedHeaders []string
	AllowedOrigins []string
}

func (c *cors) check(ctx *fasthttp.RequestCtx) error {
	if !c.Enable {
		return nil
	}

	// check origin
	origin := string(ctx.Request.Header.Peek("Origin"))
	if origin != "" && isAllowed(c.AllowedOrigins, origin) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
	}

	ctx.Response.Header.Set("Access-Control-Allow-Methods", strings.Join(c.AllowedMethods, ","))
	ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(c.AllowedHeaders, ","))

	log.Debugf("New CORs %s Request, Origin: %s, Client: %s", ctx.Method(), origin, ctx.RemoteAddr())
	return nil
}

func isAllowed(allowed []string, v string) bool {
	for _, a := range allowed {
		if a == "*" {
			return true
		}

		if strings.TrimSpace(strings.ToUpper(a)) == strings.TrimSpace(strings.ToUpper(v)) {
			return true
		}
	}

	return false
}
