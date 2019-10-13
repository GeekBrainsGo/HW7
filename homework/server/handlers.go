package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"serv/models"

	"github.com/go-chi/chi"
)

func (serv *Server) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open(filepath.Join(serv.staticDir + "/index.html"))
	data, _ := ioutil.ReadAll(file)
	if blogItems, err := models.GetBlogs(nil, serv.db); err != nil {
		serv.lg.Error("Error getting all posts", err)
	} else {
		indexTemplate := template.Must(template.New("index").Parse(string(data)))
		err := indexTemplate.ExecuteTemplate(w, "index", blogItems)
		if err != nil {
			serv.lg.WithError(err).Error("template")
		}
	}

}

func (serv *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open(filepath.Join(serv.staticDir + "/post.html"))
	data, _ := ioutil.ReadAll(file)
	postNumber := chi.URLParam(r, "id")
	indexTemplate := template.Must(template.New("index").Parse(string(data)))
	searchedPost, err := models.FindBlog(nil, serv.db, postNumber)
	if err != nil {
		serv.lg.Error("Error getting post", err)
	}
	err = indexTemplate.ExecuteTemplate(w, "index", searchedPost)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) handleGetEditPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open(filepath.Join(serv.staticDir + "/edit.html"))
	data, _ := ioutil.ReadAll(file)
	postNumber := chi.URLParam(r, "id")
	indexTemplate := template.Must(template.New("index").Parse(string(data)))
	searchedPost, err := models.FindBlog(nil, serv.db, postNumber)
	if err != nil {
		serv.lg.Error("Error getting post", err)
	}
	err = indexTemplate.ExecuteTemplate(w, "index", searchedPost)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

// handlePostEditPost - function to save edited post to DB
// @Summary saves the edited post to DB
// @Description function to update edited post in the DB
// @Accept json
// @Produce json
// @Param website body models.Blog true "Blog json struct"
// @Success 200 {object} models.Blog
// @Router /api/v1/posts/{id} [put]
func (serv *Server) HandlePostEditPost(w http.ResponseWriter, r *http.Request) {
	var post models.Blog
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		post.Update(nil, serv.db)
		resp, err := json.Marshal(post)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(resp)
		}
	}

}

// HandlePostCreatePost - function to create new post
// @Summary creates a new post
// @Description function to create a new post in the DB
// @Accept json
// @Produce json
// @Param website body models.Blog true "Blog json struct"
// @Success 200 {object} models.Blog
// @Router /api/v1/posts/ [post]
func (serv *Server) HandlePostCreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Blog
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		post.Insert(nil, serv.db)
		resp, err := json.Marshal(post)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(resp)
		}
	}
}

// HandlePostDeletePost - function to delete a post by id
// @Summary deletes a post
// @Description function to delete a post from DB
// @Accept json
// @Produce json
// @Param website body models.Blog true "Blog json struct"
// @Success 200 {object} models.Blog
// @Router /api/v1/posts/{id} [delete]
func (serv *Server) HandlePostDeletePost(w http.ResponseWriter, r *http.Request) {
	var post models.Blog
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	postNumber := chi.URLParam(r, "id")
	searchedPost, err := models.FindBlog(nil, serv.db, postNumber)
	searchedPost.Delete(nil, serv.db)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		post.Delete(nil, serv.db)
		resp, err := json.Marshal(post)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(resp)
		}
	}
}

func (serv *Server) handleSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.json")
}
