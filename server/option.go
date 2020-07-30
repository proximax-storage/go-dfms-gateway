package server

type GatewayOptions struct {
	cfg        string
	address    string
	apiAddress string
	debug      bool
}

func (opts *GatewayOptions) ApplyToConfig(cfg *Config) {
	if opts.address != "" {
		cfg.Address = opts.address
	}
	if opts.apiAddress != "" {
		cfg.ApiAddress = opts.apiAddress
	}
	if opts.apiAddress != "" {
		cfg.LogAllError = opts.debug
	}
}

type GatewayOption func(options *GatewayOptions)

func WithAddress(address string) GatewayOption {
	return func(o *GatewayOptions) {
		o.address = address
	}
}

func WithAPI(addressAPI string) GatewayOption {
	return func(o *GatewayOptions) {
		o.apiAddress = addressAPI
	}
}

func Debug(b bool) GatewayOption {
	return func(o *GatewayOptions) {
		o.debug = b
	}
}

func ConfigPath(cfg string) GatewayOption {
	return func(o *GatewayOptions) {
		if cfg != "" {
			o.cfg = resolvePath(cfg)
			return
		}
		o.cfg = resolvePath(defaultGatewayConfigPath)
	}
}

func ParseOptions(opts ...GatewayOption) *GatewayOptions {
	gopt := &GatewayOptions{}
	for _, opt := range opts {
		opt(gopt)
	}

	return gopt
}
