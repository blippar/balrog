package browser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/blippar/balrog/apk"
)

// Config defines the software configuration
type Config struct {
	Title        string            `json:"title"`
	HTTP         HTTPConfig        `json:"http"`
	Storage      StoreConfig       `json:"storage"`
	Repositories []*apk.Repository `json:"repositories"`
	Static       []*StaticFolder   `json:"static"`
}

// StaticFolder defines a staticly served storage folder
type StaticFolder struct {
	WebPrefix     string `json:"web_prefix"`
	StoragePrefix string `json:"storage_prefix"`
}

// SetDefaultConfig fills in the default values for a configuration file
func (cfg *Config) SetDefaultConfig() {
	cfg.Title = "APK Browser"
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
