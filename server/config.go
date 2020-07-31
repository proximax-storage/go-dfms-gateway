package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const defaultGatewayConfigPath = "~/.dfms/gateway_cfg.json"

func init() {
	p := resolvePath(defaultGatewayConfigPath)
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		err = saveConfig(DefaultConfig(), p)
		if err != nil {
			log.Fatalf("Save default config: ", err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Name        string
	Address     string
	ApiAddress  string
	GetOnly     bool
	LogAllError bool
}

func DefaultConfig() *Config {
	return &Config{
		Name:        "DFMS Gateway",
		Address:     ":5000",
		ApiAddress:  "http://localhost:6366",
		GetOnly:     true,
		LogAllError: false,
	}
}

func saveConfig(cfg *Config, path string) error {
	if path == "" {
		path = defaultGatewayConfigPath
	}
	path = resolvePath(path)

	content, err := json.MarshalIndent(*cfg, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, content, 0666)
	if err != nil {
		return err
	}

	return err
}

func loadConfig(path string) (*Config, error) {
	if path == "" {
		path = defaultGatewayConfigPath
	}
	path = resolvePath(path)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
