package go_dfms_gateway

type gatewayOptions struct {
	cfg     string
	address string
	debug   bool
}

func (opts *gatewayOptions) ApplyToConfig(cfg *config) {
	cfg.LogAllError = opts.debug

	if opts.address != "" {
		cfg.Address = opts.address
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
