package apk

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/0rax/apk/apk"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (s *Server) listRepoPackages(repo *apk.Repository, w http.ResponseWriter, r *http.Request) {

	arch := chi.URLParam(r, "arch")
	if arch != "" && !repo.HasArch(arch) {
		s.errorHandler(w, r, http.StatusNotFound, "specified arch is not available")
		return
	}

	pkgs, err := repo.ListPackages(r.Context(), arch)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	switch render.GetAcceptedContentType(r) {
	case render.ContentTypeJSON:
		render.JSON(w, r, pkgs)
	default:
		err = s.tmpl.ExecuteTemplate(w, "packages.html.tmpl", map[string]interface{}{
			"Title":             repo.Name,
			"Repositories":      s.Repositories,
			"CurrentArch":       arch,
			"CurrentRepository": repo,
			"Packages":          pkgs,
		})
		fmt.Println(err)
	}
}

func (s *Server) serveRepoPackages(repo *apk.Repository, w http.ResponseWriter, r *http.Request) {

	arch := chi.URLParam(r, "arch")
	file := chi.URLParam(r, "file")
	if arch != "" && !repo.HasArch(arch) {
		s.errorHandler(w, r, http.StatusNotFound, "specified arch is not available")
		return
	}

	o, err := repo.GetObject(r.Context(), path.Join(arch, file))
	if err != nil {
		if os.IsNotExist(err) {
			s.errorHandler(w, r, http.StatusNotFound, "requested package does not exists")
		} else {
			s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.Header().Add("Content-Type", "application/x-tgz")
	io.Copy(w, o)
}

func (s *Server) serveRepoKeys(repo *apk.Repository, w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "key")
	k, ok := repo.KeyName[key]
	if !ok {
		s.errorHandler(w, r, http.StatusNotFound, "specified key is not available")
		return
	}
	o, err := repo.GetObject(r.Context(), k)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	io.Copy(w, o)
}
