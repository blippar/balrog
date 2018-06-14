package browser

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Server defines an APK Repository UI & API HTTP server
type Server struct {
	*Config

	tmpl   *template.Template
	router http.Handler
}

// NewServer creates a new APK Repository server
func NewServer(cfgPath string) *Server {
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		logrus.WithError(err).Fatal("newServerFailed")
	}
	srv := &Server{
		Config: cfg,
	}
	return srv
}

// Init initlializes the required subsystem necessary for an APK Repository server based on its current configuration
func (s *Server) Init() (err error) {

	// Init storage
	if err = s.Storage.Init(); err != nil {
		return err
	}

	// Init templator
	if s.tmpl, err = template.New("").ParseGlob(path.Join(s.HTTP.Templates, "*.html.tmpl")); err != nil {
		return err
	}

	// Init router
	if err = s.initRouter(); err != nil {
		return err
	}

	return nil
}

// Run starts the HTTP Server
func (s *Server) Run() error {

	logrus.Info("serverInit")
	if err := s.Init(); err != nil {
		logrus.WithError(err).Fatal("serverInitFailed")
		return err
	}

	h := http.Server{
		Addr:    s.HTTP.Addr,
		Handler: s.router,
	}
	proto := "http"
	if h.TLSConfig != nil {
		proto = "https"
	}

	errCh := make(chan error, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() { errCh <- h.ListenAndServe() }()
	logrus.WithField("addr", fmt.Sprintf("%s://%s", proto, s.HTTP.Addr)).Info("serverRunning")

	select {
	case err := <-errCh:
		logrus.WithError(err).Error("serverStoppedWithError")
		return err
	case sig := <-sigs:
		logrus.WithField("signal", sig).Info("serverInterupt")
	}

	logrus.Info("serverStopping")
	err := h.Shutdown(context.Background())
	logrus.Info("serverStopped")

	return err
}
