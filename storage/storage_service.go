package storage

import (
	"context"
	"errors"
	"io"
)

// Service defines a standard storage service suitable for use in the APK server
type Service interface {
	Init() error
	WithPrefix(string) Service
	GetObject(context.Context, string) (io.ReadCloser, error)
}

// Err.. defines potential errors returned by a storage.Service
var (
	ErrObjectNotFound = errors.New("Object not found")
)
