package browser

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

func (s *Server) listRepositories(w http.ResponseWriter, r *http.Request) {

	switch render.GetAcceptedContentType(r) {
	case render.ContentTypeJSON:
		repo := make(map[string][]string, len(s.Repositories))
		for _, r := range s.Repositories {
			repo[r.Name] = r.Arch
		}
		render.JSON(w, r, repo)
	default:
		err := s.tmpl.ExecuteTemplate(w, "main.html.tmpl", map[string]interface{}{
			"Title":        s.Title,
			"Repositories": s.Repositories,
		})
		if err != nil {
			logrus.WithError(err).Warning("templateError")
		}
	}
}
