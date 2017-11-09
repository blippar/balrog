package apk

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
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
			"Title":        "APK Repository",
			"Repositories": s.Repositories,
		})
		fmt.Println(err)
	}
}
