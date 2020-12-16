package go_dfms_gateway

type gatewayOptions struct {
	cfg     string
	address string
	debug   bool

	enableCors     bool
	allowedMethods []string
	allowedHeaders []string
	allowedOrigins []string
}

func (opts *gatewayOptions) ApplyToConfig(cfg *config) {
	cfg.LogAllError = opts.debug
	cfg.CORs.Enable = opts.enableCors

	if opts.address != "" {
		cfg.Address = opts.address
	}

	if cfg.CORs.Enable {
		cfg.GetOnly = !opts.enableCors

		if opts.allowedMethods != nil {
			cfg.CORs.AllowedMethods = opts.allowedMethods
		}

		if opts.allowedHeaders != nil {
			cfg.CORs.AllowedHeaders = opts.allowedHeaders
		}

		if opts.allowedOrigins != nil {
			cfg.CORs.AllowedOrigins = opts.allowedOrigins
		}
	}
}

type GatewayOption func(options *gatewayOptions)

func WithAddress(address string) GatewayOption {
	return func(o *gatewayOptions) {
		o.address = address
	}
}

func Debug(b bool) GatewayOption {
	return func(o *gatewayOptions) {
		o.debug = b
	}
}

func ConfigPath(cfg string) GatewayOption {
	return func(o *gatewayOptions) {
		if cfg != "" {
			o.cfg = resolvePath(cfg)
			return
		}
		o.cfg = resolvePath(defaultGatewayConfigPath)
	}
}

func ParseOptions(opts ...GatewayOption) *gatewayOptions {
	gopt := &gatewayOptions{}

	for _, opt := range opts {
		opt(gopt)
	}

	return gopt
}

func EnableCORs(b bool) GatewayOption {
	return func(o *gatewayOptions) {
		o.enableCors = b
	}
}

func AllowedOrigins(origins ...string) GatewayOption {
	return func(o *gatewayOptions) {
		o.allowedOrigins = origins
	}
}

func AllowedHeaders(headers ...string) GatewayOption {
	return func(o *gatewayOptions) {
		o.allowedHeaders = headers
	}
}
func AllowedMethods(methods ...string) GatewayOption {
	return func(o *gatewayOptions) {
		o.allowedMethods = methods
	}
}
