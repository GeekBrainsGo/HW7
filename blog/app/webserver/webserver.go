package webserver

import (
	"blog/app/models"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// WebServer ...
type WebServer struct {
	router   *chi.Mux
	logger   *logrus.Logger
	database *mgo.Session
	config   *Config
}

func newServer(db *mgo.Session, config *Config) *WebServer { // 1
	serv := &WebServer{
		router:   chi.NewRouter(),
		logger:   logrus.New(),
		database: db,
		config:   config,
	}

	serv.configureRouter()

	logrusLevel, _ := logrus.ParseLevel(config.LogLevel)
	serv.logger.SetLevel(logrusLevel)

	return serv
}

// Start ...
func Start(config *Config) error { // 2
	db, err := newSession(config.DatabaseConnectionString)
	if err != nil {
		return err
	}

	defer db.Close()
	serv := newServer(db, config)
	return http.ListenAndServe(config.BindAddr, serv)
}

func newSession(dsnURL string) (*mgo.Session, error) { // 2

	session, err := mgo.Dial(dsnURL)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (serv *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) { // 1
	serv.router.ServeHTTP(w, r)
}

func (serv *WebServer) configureRouter() { // 1
	//routes
	serv.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	serv.router.HandleFunc("/list", serv.postListHandle())

	serv.router.HandleFunc("/view/{postID}", serv.postViewHandle())

	serv.router.HandleFunc("/delete/{postID}", serv.postDeleteHandle())

	serv.router.HandleFunc("/create", serv.postCreateHandle())

	// API routes
	serv.router.Route("/api/v1", func(router chi.Router) {
		router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(serv.config.URLShema+"://"+serv.config.Hostname+serv.config.BindAddr+serv.config.SwaggerPath)))
		router.Get("/"+serv.config.SwaggerFile, serv.HandleSwagger())

		router.Get("/posts", serv.apiPostListHandle())
		router.Get("/post/{postID}", serv.apiPostGetHandle())
		router.Put("/post", serv.apiPostCreateHandle())
	})

}

func (serv *WebServer) postListHandle() http.HandlerFunc { // 1

	type PageModel struct {
		Title string
		Data  interface{}
	}

	return func(w http.ResponseWriter, r *http.Request) { // 3

		conn := serv.database.DB("blog").C("posts")

		var posts models.PostItemsSlice

		err := conn.Find(bson.M{}).All(&posts)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		pageData := PageModel{
			Title: "Posts List",
			Data:  posts,
		}

		templ := template.Must(template.New("page").ParseFiles("./templates/blog/List.tpl", "./templates/common.tpl"))
		err = templ.ExecuteTemplate(w, "page", pageData)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

	}
}

func (serv *WebServer) postViewHandle() http.HandlerFunc { // 1

	type PageModel struct {
		Title string
		Data  interface{}
	}

	return func(w http.ResponseWriter, r *http.Request) { // 3

		postID := chi.URLParam(r, "postID")

		conn := serv.database.DB("blog").C("posts")

		var post models.Post

		err := conn.Find(bson.M{"id": postID}).One(&post)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		pageData := PageModel{
			Title: "View Post",
			Data:  post,
		}

		templ := template.Must(template.New("page").ParseFiles("./templates/blog/View.tpl", "./templates/common.tpl"))
		err = templ.ExecuteTemplate(w, "page", pageData)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

	}
}

func (serv *WebServer) postDeleteHandle() http.HandlerFunc { // 1

	return func(w http.ResponseWriter, r *http.Request) { // 2

		postID := chi.URLParam(r, "postID")

		conn := serv.database.DB("blog").C("posts")

		err := conn.Remove(bson.M{"id": postID})
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return

	}
}

func (serv *WebServer) postCreateHandle() http.HandlerFunc { // 1

	return func(w http.ResponseWriter, r *http.Request) { // 2

		var newPost models.Post

		newPost.ID = uuid.NewV4().String()
		newPost.Title = "New Post Title"
		newPost.Short = "Short body"
		newPost.Body = "Content body"

		conn := serv.database.DB("blog").C("posts")

		err := conn.Insert(newPost)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return

	}
}

// Posts list - All posts
// @Description Returns all posts
// @Tags system
// @Success 200 {string} string
// @Router /api/v1/posts [get]
func (serv *WebServer) apiPostListHandle() http.HandlerFunc { // 1

	return func(w http.ResponseWriter, r *http.Request) { // 2

		conn := serv.database.DB("blog").C("posts")

		var posts models.PostItemsSlice

		err := conn.Find(bson.M{}).All(&posts)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		serv.respondJSON(w, r, http.StatusOK, posts)
	}
}

// Post get - get one post
// @Description Returns one post
// @Param id path string true "Example: 1054497f-7c0b-4579-b4f2-524f58c712f7"
// @Tags system
// @Success 200 {string} string
// @Router /api/v1/post/{id} [get]
func (serv *WebServer) apiPostGetHandle() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		postID := chi.URLParam(r, "postID")

		conn := serv.database.DB("blog").C("posts")

		var post models.Post

		err := conn.Find(bson.M{"id": postID}).One(&post)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		serv.respondJSON(w, r, http.StatusOK, post)
	}
}

// Post create - create new post
// @Description Returns one post
// @Tags system
// @Success 201 {string} string
// @Router /api/v1/post [put]
func (serv *WebServer) apiPostCreateHandle() http.HandlerFunc { // 1

	return func(w http.ResponseWriter, r *http.Request) { // 5

		newPost := models.Post{}

		jsonData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			serv.errorAPI(w, r, http.StatusBadRequest, err)
			return
		}

		err = json.Unmarshal(jsonData, &newPost)
		if err != nil {
			serv.errorAPI(w, r, http.StatusBadRequest, err)
			return
		}

		if newPost.ID == "" {
			newPost.ID = uuid.NewV4().String()
		}

		conn := serv.database.DB("blog").C("posts")

		err = conn.Insert(newPost)
		if err != nil {
			serv.errorAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		location := "/api/v1/post/" + newPost.ID

		serv.respondJSON(w, r, http.StatusCreated, location)
	}
}

// HandleSwagger - Returns swagger.json docs
// @Description Returns swagger.json docs
// @Tags system
// @Success 200 {string} string
// @Router /api/v1/docs/swagger.json [get]
func (serv *WebServer) HandleSwagger() http.HandlerFunc { // 1
	return func(w http.ResponseWriter, r *http.Request) { // 1
		http.ServeFile(w, r, serv.config.SwaggerFile)
	}
}

func (serv *WebServer) errorAPI(w http.ResponseWriter, r *http.Request, code int, err error) { // 1
	serv.respondJSON(w, r, code, map[string]string{"error": err.Error()})
}

func (serv *WebServer) respondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) { // 2
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (serv *WebServer) respondWhithTemplate(w http.ResponseWriter, r *http.Request, code int, data interface{}) { // 2
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
