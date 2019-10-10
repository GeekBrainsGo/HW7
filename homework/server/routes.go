package server

import "github.com/go-chi/chi"

func (serv *Server) bindRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/edit/{id}", serv.handlePostEditPost)
			r.Post("/create", serv.handlePostCreatePost)
			r.Post("/delete/{id}", serv.handlePostDeletePost)
		})
		r.Get("/", serv.handleGetIndex)
		r.Get("/post/{id}", serv.handleGetPost)
		r.Get("/edit/{id}", serv.handleGetEditPost)
	})
}
