package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"serv/models"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	// TEMPLATESDIR строка описвывающая раздел шаблонов страниц
	TEMPLATESDIR = "/templates"
)

// Server - объект сервера
type Server struct {
	lg *logrus.Logger
	// db         *sql.DB
	db         *gorm.DB
	rootDir    string
	Page       models.Page
	dictionary map[string]*template.Template
}

// New - создаёт новый экземпляр сервера
// func New(lg *logrus.Logger, rootDir string, db *sql.DB) *Server {
func New(lg *logrus.Logger, rootDir string, db *gorm.DB) *Server {

	// initial
	fileBase := path.Join(rootDir, TEMPLATESDIR, "base.html")
	fileMenu := path.Join(rootDir, TEMPLATESDIR, "menu.html")
	tMap := make(map[string]*template.Template)

	fileAdd := path.Join(rootDir, TEMPLATESDIR, "blogs.html")
	temp := template.Must(template.ParseFiles(fileBase, fileMenu, fileAdd))
	tMap["BLOGS"] = temp

	fileAdd = path.Join(rootDir, TEMPLATESDIR, "blog.html")
	temp = template.Must(template.ParseFiles(fileBase, fileMenu, fileAdd))
	tMap["BLOG"] = temp

	return &Server{
		lg:         lg,
		db:         db,
		rootDir:    rootDir,
		Page:       models.Page{},
		dictionary: tMap,
	}
}

// Start - запускает сервер
func (serv *Server) Start(addr string) error {
	r := chi.NewRouter()
	serv.bindRoutes(r)
	serv.lg.Debug("server is started ...")
	return http.ListenAndServe(addr, r)
}

func (serv *Server) bindRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.getTemplateHandler)
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/docs/swagger.json", serv.HandleSwaggerJSON)
			r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/api/v1/docs/swagger.json")))
			r.Put("/blog/{id}", serv.putBlogHandler)
			r.Get("/blog/{id}", serv.editBlogHandler)
			r.Post("/blog", serv.addBlogHandler)
			r.Get("/blog", serv.addGetBlogHandler)
			r.Delete("/blog/{id}", serv.deleteBlogHandler)
		})
	})
}

// SendErr - отправляет ошибку пользователю и логирует её
func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("server error")
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

// SendInternalErr - отправляет 500 ошибку
func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}
