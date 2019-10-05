package main

import (
	"HW7-master/models"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	STATIC_DIR = "./www/"
)

func main() {
	lg := logrus.New()
	ctx := context.Background()
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	_ = client.Connect(ctx)
	db := client.Database("blog")

	serv := NewServer(lg, db, ctx)
	serv.Start(":8080")

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan

	serv.Stop()
}

type Server struct {
	lg    *logrus.Logger
	db    *mongo.Database
	ctx   context.Context
	mux      *chi.Mux
	server   *http.Server
	swagURL  string
	swagPath string
	Title string
	Posts models.Posts
}

// NewServer - Create a new server instance
func NewServer(lg *logrus.Logger, db *mongo.Database, ctx context.Context) *Server {
	serv := &Server{
		lg:  lg,
		db:  db,
		ctx: ctx,
		mux: chi.NewMux(),

		swagURL: "http://localhost:8080/api/v1/docs/swagger.json",
		swagPath: "docs/swagger.json",

		Title: "BLOG",
		Posts: models.Posts {},
	}

	return serv
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

func (serv *Server) configureRoutes() {
	serv.mux.Route("/", func(r chi.Router) {
		fileServer := http.FileServer(http.Dir(STATIC_DIR))
		r.Handle("/static/*", fileServer)

		r.Route("/", func(r chi.Router) {
			r.Get("/", serv.HandleGetIndex)
			r.Get("/post/{id}", serv.HandleGetPost)
			r.Get("/post/create", serv.HandleGetEditPost)
			r.Get("/post/{id}/edit", serv.HandleGetEditPost)
		})

		r.Route("/api/v1/", func(r chi.Router) {
			r.Post("/post", serv.HandleEditPost)
			r.Put("/post/{id}", serv.HandleEditPost)
		})
	})
}

// HandleSwagger - Returns swagger.json docs
// @Description Returns swagger.json docs
// @Tags system
// @Success 200 {string} string
// @Router /docs/swagger.json [get]
func (serv *Server) HandleSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, serv.swagPath)
}

func (serv *Server) AddOrUpdatePost(newPost models.Post) (models.Post) {
	updPost, err := newPost.Update(serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("AddOrUpdatePost")
		newPost, _ := newPost.Create(serv.ctx, serv.db)
		return *newPost
	}

	return *updPost
}

func (serv *Server) HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/index.gohtml")
	data, _ := ioutil.ReadAll(file)

	posts, err := models.GetPosts(serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetIndex")
		posts = models.Posts{}
	}

	serv.Posts = posts

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetIndexTemplate")
	}
}

func (serv *Server) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/post.gohtml")
	data, _ := ioutil.ReadAll(file)

	postIDStr := chi.URLParam(r, "id")
	postID, _ := primitive.ObjectIDFromHex(postIDStr)

	post, err := models.GetPost(postID, serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetPost")
		post = &models.Post{}
	}

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetPostTemplate")
	}
}

func (serv *Server) HandleGetEditPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/templates/edit_post.gohtml")
	data, _ := ioutil.ReadAll(file)

	postIDStr := chi.URLParam(r, "id")
	postID, _ := primitive.ObjectIDFromHex(postIDStr)

	post, err := models.GetPost(postID, serv.ctx, serv.db)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetEditPost")
		post = &models.Post{}
	}

	templ := template.Must(template.New("page").Parse(string(data)))
	err = templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("HandleGetEditPostTemplate")
	}
}

// HandleEditPost - Create and Edit post
// @Description Create and Edit post
// @Tags create-post, edit-post
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 400
// @Router /api/v1/post [post]
// @Router /api/v1/post/{id} [put]
func (serv *Server) HandleEditPost(w http.ResponseWriter, r *http.Request) {

	/*
	{"id":"5d986adf81f8d848d93e11ca", "Title":"Пост 4", "Text":"Новый очень интересный текст", "Labels":["l1","l2"]}
	*/

	decoder := json.NewDecoder(r.Body)
	var inPostItem models.Post

	err := decoder.Decode(&inPostItem)
	if err != nil {
		serv.lg.WithError(err).Error("HandleEditPost")
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	newPost := serv.AddOrUpdatePost(inPostItem)
	respondWithJSON(w, http.StatusOK, newPost)
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

