package server

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (serv *Server) bindRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/api/v1/docs/swagger.json")))
			r.Get("/docs/swagger.json", serv.handleSwaggerJSON)
			r.Route("/posts", func(r chi.Router) {
				r.Put("/", serv.HandlePostEditPost)
				r.Post("/", serv.HandlePostCreatePost)
				r.Delete("/{id}", serv.HandlePostDeletePost)
			})
		})
		r.Get("/", serv.handleGetIndex)
		r.Get("/post/{id}", serv.handleGetPost)
		r.Get("/edit/{id}", serv.handleGetEditPost)
	})
}
