package browser

import (
	"fmt"

	"github.com/blippar/balrog/storage"
	"github.com/blippar/balrog/storage/afero"
	"github.com/blippar/balrog/storage/minio"
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
