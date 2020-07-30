package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const defaultGatewayConfigPath = "~/.dfms/gateway_cfg.json"

func init() {
	_, err := os.Stat(resolvePath(defaultGatewayConfigPath))
	if os.IsNotExist(err) {
		err = saveConfig(DefaultConfig())
		if err != nil {
			log.Fatal("Save default config: ", err)
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

func saveConfig(cfg *Config) error {
	content, err := json.MarshalIndent(*cfg, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(resolvePath(defaultGatewayConfigPath), content, 0666)
	if err != nil {
		return err
	}

	return err
}

func loadConfig(path string) (*Config, error) {
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
