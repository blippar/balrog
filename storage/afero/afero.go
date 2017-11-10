package afero

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	"github.com/blippar/balrog/storage"
	"github.com/spf13/afero"
)

var _ = storage.Service(&Storage{})

// Storage defines a storage service based on your filesystem
type Storage struct {
	afero.Fs `json:"-"`

	Prefix string `json:"prefix"`
	Cache  int64  `json:"cache"` // in seconds
}

// Init configures the underlying storage engine
func (s *Storage) Init() error {
	s.Fs = afero.NewOsFs()
	if s.Cache > 0 {
		lfs := afero.NewMemMapFs()
		s.Fs = afero.NewCacheOnReadFs(s.Fs, lfs, time.Duration(s.Cache)*time.Second)
	}
	return nil
}

// WithPrefix returns a new storage service with the new prefixed appended to the current one
func (s *Storage) WithPrefix(prefix string) storage.Service {
	srv := *s
	srv.Prefix = path.Join(s.Prefix, prefix)
	return &srv
}

// GetObject returns the specified object if found
func (s *Storage) GetObject(_ context.Context, n string) (io.ReadCloser, error) {
	f, err := s.Open(path.Join(s.Prefix, n))
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	} else if info.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}
