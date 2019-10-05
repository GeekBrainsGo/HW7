// Package server implements blog sever with mongodb and swagger.
package server

import (
	"HW7/swaggerblog/models"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// Server stands for server struct.
type Server struct {
	rootDir       string
	templatesDir  string
	indexTemplate string
	swagURL       string
	swagPath      string
	ctx           context.Context
	db            *mongo.Database
	lg            *logrus.Logger
	mux           *chi.Mux
	Page          models.Page
	server        *http.Server
}

// New creates new server.
func New(ctx context.Context, lg *logrus.Logger, rootDir string, db *mongo.Database) *Server {
	return &Server{
		lg:            lg,
		db:            db,
		rootDir:       rootDir,
		templatesDir:  "/templates",
		indexTemplate: "index.html",
		ctx:           ctx,
		mux:           chi.NewMux(),
	}
}

// SetSwagger sets swagger.json path.
func (s *Server) SetSwagger(swagPath, swagURL string) {
	s.swagPath = swagPath
	s.swagURL = swagURL
}

// Start starts new server.
func (s *Server) Start(addr string) *Server {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	s.routes()

	go func() {
		s.lg.Info("server started")
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			s.lg.Errorf("server: %s", err)
		}
	}()

	return s
}

// Stop stops the server.
func (s *Server) Stop() error {
	s.lg.Info("stopping server")
	return s.server.Shutdown(s.ctx)
}

// SendErr sends and log error to user.
func (s *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	s.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := models.ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr sends 500 error.
func (s *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	s.SendErr(w, err, http.StatusInternalServerError, obj)
}
