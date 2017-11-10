package browser

import (
	"net/http"

	"github.com/go-chi/render"
)

type httpError struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Cause  string `json:"cause,omitempty"`
}

func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, code int, cause string) {

	var jErr = httpError{
		Code:   code,
		Cause:  cause,
		Status: http.StatusText(code),
	}
	w.WriteHeader(code)

	switch render.GetAcceptedContentType(r) {
	case render.ContentTypeJSON:
		render.JSON(w, r, jErr)
	default:
		s.tmpl.ExecuteTemplate(w, "errors.html.tmpl", map[string]interface{}{
			"Title":        s.Title,
			"Repositories": s.Repositories,
			"Error":        jErr,
		})
	}

}
