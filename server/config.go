package server

type GatewayConfig struct {
	Name       string
	Address    string
	AddressAPI string
	GetOnly    bool
	LogError   bool
}

func DefaultGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Name:       "DFMS Gateway",
		Address:    ":5000",
		AddressAPI: "http://localhost:6366",
		GetOnly:    true,
		LogError:   true,
	}
}
