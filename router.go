package browser

import (
	"net/http"
	"os"
	"path"

	"github.com/blippar/balrog/apk"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) initRouter() error {

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RealIP)
	r.Use(s.AccessLogMiddleware)

	// Error handlers
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.errorHandler(w, r, http.StatusNotFound, "")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		s.errorHandler(w, r, http.StatusMethodNotAllowed, "")
	})

	// Router
	r.Get("/", s.listRepositories)
	for _, rep := range s.Repositories {
		rep.Init(s.Storage.Service)
		r.Mount("/"+rep.Name+"/", s.newRepositoryRouter(rep))
	}
	for _, st := range s.Static {
		r.Mount("/"+st.WebPrefix+"/", s.newStaticRouter(st))
	}
	r.Mount("/dist/", s.newDistRouter("/dist", s.HTTP.Dist))

	// Attach router to Server
	s.router = r

	return nil
}

func (s *Server) newDistRouter(prefix, folder string) http.Handler {

	r := chi.NewRouter()

	if !path.IsAbs(folder) {
		workDir, _ := os.Getwd()
		folder = path.Join(workDir, folder)
	}

	fs := http.StripPrefix(prefix, http.FileServer(http.Dir(folder)))
	r.Get("/*", fs.ServeHTTP)
	return r
}

// newRepositoryRouter creates a new subrouter for the passed reository
func (s *Server) newRepositoryRouter(repo *apk.Repository) http.Handler {
	r := chi.NewRouter()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.listRepoPackages(repo, w, r)
	}))
	if repo.KeyName != nil {
		r.Get("/{key:.+\\.rsa\\.pub}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.serveRepoKeys(repo, w, r)
		}))
	}
	r.Get("/{arch}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.listRepoPackages(repo, w, r)
	}))
	r.Get("/{arch}/{file}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.serveRepoPackages(repo, w, r)
	}))
	return r
}

// newStaticRouter creates a new subrouter for the passed static folder
func (s *Server) newStaticRouter(st *StaticFolder) http.Handler {
	r := chi.NewRouter()
	store := s.Storage.WithPrefix(st.StoragePrefix)
	r.Get("/{file}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.serveStaticFile(store, w, r)
	}))
	return r
}
