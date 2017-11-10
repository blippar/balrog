package apk

import (
	"context"
	"os"
	"path"

	"github.com/blippar/balrog/storage"
)

// Repository defines an APK repository configuration and underlying storage
type Repository struct {
	storage.Service `json:"-"`

	Name   string   `json:"name"`
	Arch   []string `json:"available_arch"`
	Prefix string   `json:"storage_prefix"`
	Keys   []string `json:"public_keys"`

	KeyName map[string]string `json:"-"`
}

// Init configure a Repository using the passed storage.Service and its current configuration
func (r *Repository) Init(store storage.Service) error {
	r.Service = store.WithPrefix(r.Prefix)
	r.KeyName = make(map[string]string, len(r.Keys))
	for _, k := range r.Keys {
		name := path.Base(k)
		r.KeyName[name] = k
	}
	return nil
}

// ListPackages returns the list of available package in the repository based on the passed arch
func (r *Repository) ListPackages(ctx context.Context, arch string) ([]Package, error) {

	pkgs := make([]Package, 0)

	if arch != "" {
		o, err := r.WithPrefix(arch).GetObject(ctx, IndexTar)
		if os.IsNotExist(err) {
			return pkgs, nil
		} else if err != nil {
			return pkgs, err
		}
		defer o.Close()
		return ListPackageFromIndex(o)
	}

	for _, a := range r.Arch {
		ipkgs, err := r.ListPackages(ctx, a)
		if err != nil && !os.IsNotExist(err) {
			return pkgs, err
		}
		pkgs = append(pkgs, ipkgs...)
	}
	return pkgs, nil
}

// HasArch allows verification that a specific arch is available in this repository
func (r *Repository) HasArch(arch string) bool {
	for _, a := range r.Arch {
		if a == arch {
			return true
		}
	}
	return false
}
