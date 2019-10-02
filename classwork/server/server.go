package server

import (
	"context"
	"encoding/json"
	"net/http"
	"serv/models"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// Server - Server instance
type Server struct {
	swagURL  string
	swagPath string
	lg       *logrus.Logger
	ctx      context.Context
	mux      *chi.Mux
	server   *http.Server
}

// NewServer - Create a new server instance
func NewServer(ctx context.Context, lg *logrus.Logger) *Server {
	serv := &Server{
		lg:  lg,
		ctx: ctx,
		mux: chi.NewMux(),
	}
	return serv
}

// SetSwagger - Setting swagger.json path
func (serv *Server) SetSwagger(swagPath, swagURL string) {
	serv.swagPath = swagPath
	serv.swagURL = swagURL
}

// Start - Starts the server
func (serv *Server) Start(addr string) *Server {
	serv.server = &http.Server{
		Addr:    addr,
		Handler: serv.mux,
	}
	serv.configureRoutes()

	go func() {
		serv.lg.Info("starting server")
		if err := serv.server.ListenAndServe(); err != http.ErrServerClosed {
			serv.lg.Errorf("server: %s", err)
		}
	}()

	return serv
}

// Stop - Stops the server
func (serv *Server) Stop() error {
	serv.lg.Info("stopping server")
	return serv.server.Shutdown(serv.ctx)
}

// SendErr - Sends error
func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := models.ServErr{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr - Sends internal server error
func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}

// SendJSON - Wrapper for structure sending
func (serv *Server) SendJSON(w http.ResponseWriter, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}
