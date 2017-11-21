package apk

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// ABuild default values for packages and indexes
const (
	IndexTar   = "APKINDEX.tar.gz"
	IndexFile  = "APKINDEX"
	PackageExt = ".apk"
	PubKeyExt  = ".rsa.pub"
)

// Index represents the content of the APKINDEX file in an index tar
type Index struct {
	Packages []Package
}

// Package defines a package entry in an APKINDEX file
type Package struct {
	C             string    `apk:"C" json:"-"`
	Package       string    `apk:"P" json:"name"`
	Version       string    `apk:"V" json:"version"`
	Description   string    `apk:"T" json:"description"`
	Arch          string    `apk:"A" json:"arch"`
	PackageSize   string    `apk:"S" json:"size"`
	InstalledSize string    `apk:"I" json:"installed_size"`
	URL           string    `apk:"U" json:"url"`
	License       string    `apk:"L" json:"license"`
	Origin        string    `apk:"o" json:"origin"`
	Maintainer    string    `apk:"m" json:"maintainer"`
	BuildTime     time.Time `apk:"t" json:"build_time"`
	Commit        string    `apk:"c" json:"build_commit"`
	Dependencies  []string  `apk:"D" json:"dependencies"`
	Provides      []string  `apk:"p" json:"provides"`
	APK           string    `apk:"-" json:"apk_url"`
}

func parseIndex(index io.Reader) (*Index, error) {
	idx, err := ioutil.ReadAll(index)
	if err != nil {
		return nil, err
	}
	pkgs := make([]Package, 0)
	for _, pkg := range bytes.Split(idx, []byte("\n\n")) {
		p := Package{}
		for _, l := range bytes.Split(pkg, []byte("\n")) {
			a := bytes.SplitN(l, []byte{':'}, 2)
			if len(a) < 2 {
				continue
			}
			switch k, v := string(a[0]), string(a[1]); k {
			case "C":
				p.C = v
			case "P":
				p.Package = v
			case "V":
				p.Version = v
			case "A":
				p.Arch = v
			case "S":
				p.PackageSize = v
			case "I":
				p.InstalledSize = v
			case "T":
				p.Description = v
			case "U":
				p.URL = v
			case "L":
				p.License = v
			case "o":
				p.Origin = v
			case "m":
				p.Maintainer = v
			case "t":
				d, _ := strconv.ParseInt(v, 10, 64)
				p.BuildTime = time.Unix(d, 0)
			case "c":
				p.Commit = v
			case "D":
				p.Dependencies = strings.Split(v, " ")
			case "p":
				p.Provides = strings.Split(v, " ")
			}
		}
		if p.Package != "" {
			p.APK = fmt.Sprintf("%s-%s.apk", p.Package, p.Version)
			pkgs = append(pkgs, p)
		}
	}
	return &Index{Packages: pkgs}, nil
}

// ListPackageFromIndex reads an APKINDEX.tar.gz and return a list of packages
func ListPackageFromIndex(o io.Reader) ([]Package, error) {
	gr, err := gzip.NewReader(o)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(gr)
	var index io.Reader
	for {
		header, nErr := tarReader.Next()
		if nErr == io.EOF {
			break
		}
		if nErr != nil {
			return nil, nErr
		}
		if header.Typeflag == tar.TypeReg && header.Name == IndexFile {
			index = tarReader
			break
		}
	}
	if index == nil {
		return nil, fmt.Errorf("Can't retrieve %s", IndexFile)
	}
	idx, err := parseIndex(index)
	if err != nil {
		return nil, err
	}
	return idx.Packages, nil
}
