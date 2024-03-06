package rest

import (
	"context"
	"net/http"
	"rest/database"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

var log = logrus.New()

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler: handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	database.GetDB()


	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
