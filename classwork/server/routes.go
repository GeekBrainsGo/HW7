package server

import (
	"net/http"

	// Swagger docs for this server
	_ "serv/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi"
)

func (serv *Server) configureRoutes() {
	serv.mux.Route("/", func(r chi.Router) {
		// Static folder handler
		r.Handle("/*", http.FileServer(http.Dir("./assets")))

		// API
		r.Route("/api/v1", func(r chi.Router) {
			// Serv swagger.json and swagger local page
			r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(serv.swagURL)))
			r.Get("/"+serv.swagPath, serv.HandleSwagger)

			// Echo test
			r.Post("/echo", serv.HandleEcho)
		})
	})
}
