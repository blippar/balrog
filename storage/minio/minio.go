package minio

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/0rax/apk/storage"
	"github.com/minio/minio-go"
)

const minioErrNoSuchKey = "NoSuchKey"

// Storage defines a storage service based on minio
type Storage struct {
	*minio.Client `json:"-"`

	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"access_key"`
	PrivateKey string `json:"private_key"`
	SSL        bool   `json:"ssl"`
	Bucket     string `json:"bucket"`
	Prefix     string `json:"prefix"`
}

// Init configures the underlying storage engine
func (s *Storage) Init() (err error) {
	s.Client, err = minio.New(s.Endpoint, s.AccessKey, s.PrivateKey, s.SSL)
	if err != nil {
		return err
	}
	if ok, err := s.Client.BucketExists(s.Bucket); !ok {
		return fmt.Errorf("Bucket %s does not exists", s.Bucket)
	} else if err != nil {
		return err
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
func (s Storage) GetObject(ctx context.Context, filePath string) (io.ReadCloser, error) {

	fmt.Println(path.Join(s.Prefix, filePath))

	o, err := s.Client.GetObjectWithContext(ctx, s.Bucket, path.Join(s.Prefix, filePath), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	_, err = o.Stat()
	if err != nil {
		e := minio.ToErrorResponse(err)
		if e.Code == minioErrNoSuchKey {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	return o, nil
}
