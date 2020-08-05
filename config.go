package go_dfms_gateway

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const defaultGatewayConfigPath = "~/.dfms_gateway/config.json"

func init() {
	p := resolvePath(defaultGatewayConfigPath)
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		err = saveConfig(defaultConfig(), p)
		if err != nil {
			log.Fatalf("Save default config: ", err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
}

type config struct {
	Name        string
	Address     string
	GetOnly     bool
	LogAllError bool
}

func defaultConfig() *config {
	return &config{
		Name:        "DFMS Gateway",
		Address:     ":5000",
		GetOnly:     true,
		LogAllError: false,
	}
}

func saveConfig(cfg *config, p string) error {
	if p == "" {
		p = defaultGatewayConfigPath
	}
	p = resolvePath(p)

	content, err := json.MarshalIndent(*cfg, "", "\t")
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(p), os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p, content, 0666)
	if err != nil {
		return err
	}

	return err
}

func loadConfig(path string) (*config, error) {
	if path == "" {
		path = defaultGatewayConfigPath
	}
	path = resolvePath(path)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
