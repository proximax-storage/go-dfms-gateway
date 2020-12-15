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
		Enable:         true,
		AllowedMethods: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
	}
}

type cors struct {
	Enable bool

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

	//ctx.Response.Header.Set("Access-Control-Allow-Methods", "Content-type")
	//ctx.Response.Header.Set("Access-Control-Allow-Headers", "GET")

	//check method
	requestMethod := string(ctx.Request.Header.Peek("Access-Control-Request-Method"))
	if requestMethod != "" || isAllowed(c.AllowedMethods, requestMethod) {
		ctx.Response.Header.Set("Access-Control-Allow-Methods", strings.Join(c.AllowedMethods, ","))
	}

	// check header
	requestHeaders := string(ctx.Request.Header.Peek("Access-Control-Request-Headers"))
	if requestHeaders != "" {
		headers := strings.Split(requestHeaders, ",")

		allowed := true
		for _, h := range headers {
			if !isAllowed(c.AllowedHeaders, h) {
				allowed = false
				break
			}
		}

		if allowed {
			ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(c.AllowedHeaders, ","))
		}
	}

	log.Debugf("Enable CORs for: Client: %s, %s Request, URL: %s ", ctx.RemoteAddr(), ctx.Method(), ctx.URI())
	return nil
}

func isAllowed(allowed []string, v string) bool {
	for _, a := range allowed {
		if a == "*" {
			return true
		}

		//for _, v := range values {
		if strings.TrimSpace(strings.ToUpper(a)) == strings.TrimSpace(strings.ToUpper(v)) {
			return true
		}
		//}
	}

	return false
}
