package apk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/0rax/apk/apk"
	"github.com/0rax/apk/storage"
	"github.com/0rax/apk/storage/afero"
	"github.com/0rax/apk/storage/minio"
)

// Config defines the software configuration
type Config struct {
	HTTP         HTTPConfig        `json:"http"`
	Storage      StoreConfig       `json:"storage"`
	Repositories []*apk.Repository `json:"repositories"`
}

// SetDefaultConfig fills in the default values for a configuration file
func (cfg *Config) SetDefaultConfig() {
	cfg.HTTP.Addr = ":8000"
}

// HTTPConfig defines a standard HTTP server configuration
type HTTPConfig struct {
	Addr string `json:"addr"`
}

// StoreConfig defines a configurable StorageService
type StoreConfig struct {
	storage.Service
	Type       string         `json:"type"`
	Minio      *minio.Storage `json:"minio"`
	Filesystem *afero.Storage `json:"filesystem"`
}

// Init creates the underlying storage service based on configuration
func (s *StoreConfig) Init() error {
	switch s.Type {
	case "minio":
		s.Service = s.Minio
	case "filesystem":
		s.Service = s.Filesystem
	}
	if s.Service != nil {
		return s.Service.Init()
	}
	return fmt.Errorf("config/storage: found no suitable storage for '%s'", s.Type)
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
