package server

import (
	"HW7/swaggerblog/models"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi"
)

// TemplateHandle returns template.
// @Summary serve templated page
// @Description return templated page if exist
// @Param template path string true "in form like index.html"
// @Success 200 {string} string
// @Failure 500 {string} models.ErrorModel
// @Produce html
// @Router /{template} [get]
func (s *Server) TemplateHandle(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")

	if templateName == "" {
		templateName = s.indexTemplate
	}

	file, err := os.Open(path.Join(s.rootDir, s.templatesDir, templateName))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.SendInternalErr(w, err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("").Parse(string(data))
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}

	posts, err := models.AllPosts(s.db)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}

	s.Page.Posts = posts

	if err := templ.Execute(w, s.Page); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}

// PostHandle create new post.
// @Summary create new post
// @Description Create new post and return json
// @Success 200 {string} models.Post
// @Failure 500 {string} models.ErrorModel
// @Produce json
// @Param post body models.Post true "title, author, content"
// @Router /api/v1/posts [post]
func (s *Server) PostHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	post := &models.Post{}
	if err = json.NewDecoder(r.Body).Decode(post); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	err = post.Insert(s.db)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// DeleteHandle deletes a post.
// @Summary deletes a post by id
// @Description deletes a post by id
// @Success 200 {string} string
// @Param id path string true "id of post"
// @Router /api/v1/posts/{id} [delete]
func (s *Server) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	hex := chi.URLParam(r, "id")
	post := models.Post{Hex: hex}
	if _, err := post.Delete(s.db); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PutHandle renew post.
// @Summary renew post by id
// @Description renew post by id
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Param id path string true "id of post"
// @Param post body models.Post true "title, author, content"
// @Router /api/v1/posts/{id} [put]
func (s *Server) PutHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	post := &models.Post{}
	if err = json.NewDecoder(r.Body).Decode(post); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	if _, err = post.Update(s.db); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// SwaggerHandle serve path to json for swagger.
// @Summary returns json docs
// @Description Returns swagger's json docs
// @Produce  json
// @Success 200 {string} string "ok"
// @Router /api/v1/docs/swagger.json [get]
func (s *Server) SwaggerHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.swagPath)
}
