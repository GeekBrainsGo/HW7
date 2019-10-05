package server

import (
	// _ "go_basics/packages/swaggerblog/docs"

	"github.com/go-chi/chi"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) routes() {
	s.mux.Route("/", func(r chi.Router) {
		r.Get("/{template}", s.TemplateHandle)
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(s.swagURL)))
			r.Get("/"+s.swagPath, s.SwaggerHandle)
			r.Post("/posts", s.PostHandle)
			r.Delete("/posts/{id}", s.DeleteHandle)
			r.Put("/posts/{id}", s.PutHandle)
		})
	})
}
