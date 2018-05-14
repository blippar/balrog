package storage

import (
	"context"
	"io"
)

// Service defines a standard storage service suitable for use in the APK server
type Service interface {
	Init() error
	WithPrefix(string) Service
	GetObject(context.Context, string) (io.ReadCloser, error)
}
