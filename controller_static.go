package browser

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/blippar/balrog/storage"
	"github.com/go-chi/chi"
)

func (s *Server) serveStaticFile(store storage.Service, w http.ResponseWriter, r *http.Request) {

	file := chi.URLParam(r, "file")
	o, err := store.GetObject(r.Context(), file)
	if err != nil {
		fmt.Println(err)
		if err == os.ErrNotExist {
			s.errorHandler(w, r, http.StatusNotFound, err.Error())
		} else {
			s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	io.Copy(w, o)
}
