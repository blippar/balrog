package browser

import (
	"fmt"

	"github.com/blippar/alpine-package-browser/storage"
	"github.com/blippar/alpine-package-browser/storage/afero"
	"github.com/blippar/alpine-package-browser/storage/minio"
)

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
