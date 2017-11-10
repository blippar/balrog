package browser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/blippar/alpine-package-browser/apk"
)

// Config defines the software configuration
type Config struct {
	HTTP         HTTPConfig        `json:"http"`
	Storage      StoreConfig       `json:"storage"`
	Repositories []*apk.Repository `json:"repositories"`
}

// SetDefaultConfig fills in the default values for a configuration file
func (cfg *Config) SetDefaultConfig() {
	cfg.HTTP.SetDefaultConfig()
}

func loadConfig(path string) (*Config, error) {

	cfg := &Config{}
	cfg.SetDefaultConfig()

	cfile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cfile, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
