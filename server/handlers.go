package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"serv/models"
	"strconv"

	"github.com/go-chi/chi"
)

// HandleSwaggerJSON - Возвращает swagger.json
// @Summary Возвращает swagger.json
// @Tags blog, list
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /api/v1/swagger.json [get]
func (serv *Server) HandleSwaggerJSON(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "./docs/swagger.json")
	data, _ := ioutil.ReadAll(r.Body)
	w.Write(data)

}

// getTemplateHandler - возвращает список блогов
// @Summary Список блогов
// @Description получить список блогов
// @Router / [get]
func (serv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {

	blogs, err := models.GetAllBlogs(serv.db)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	serv.Page.Title = "Мой блог"
	serv.Page.Data = blogs
	serv.Page.Command = "index"

	if err := serv.dictionary["BLOGS"].ExecuteTemplate(w, "base", serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

func getBlogIDFromRequest(r *http.Request) (uint, error) {

	bl := chi.URLParam(r, "id")
	// blogID, _ := strconv.ParseInt(bl, 10, 64)
	tempInt, err := strconv.Atoi(bl)
	if err != nil {
		return 0, err
	}
	blogID := uint(tempInt)

	return blogID, err
}

// editBlogHandler - редактирование блога
// @Summary Редактирование блога
// @Description редактирование блога
// @Tags blog, edit
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /api/v1/blog/{id} [get]
func (serv *Server) editBlogHandler(w http.ResponseWriter, r *http.Request) {

	blogID, err := getBlogIDFromRequest(r)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	blog, err := models.GetBlog(serv.db, blogID)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	serv.Page.Title = "Редактирование"
	serv.Page.Data = blog
	serv.Page.Command = "edit"

	if err := serv.dictionary["BLOG"].ExecuteTemplate(w, "base", serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// putBlogHandler - обновляет блог
// @Summary Обновить блог
// @Description обновление блога
// @Tags blog, edit
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /api/v1/blog/{id} [put]
func (serv *Server) putBlogHandler(w http.ResponseWriter, r *http.Request) {

	blogID, err := getBlogIDFromRequest(r)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)

	blog := models.Blog{}
	_ = json.Unmarshal(data, &blog)

	blog.ID = blogID

	if err := blog.UpdateBlog(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	http.Redirect(w, r, "/", 301)
}

// addGetBlogHandler - добавление блога
// @Summary Добавление нового блога
// @Description добавление нового блога
// @Tags blog, add
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /api/v1/blog/{id} [get]
func (serv *Server) addGetBlogHandler(w http.ResponseWriter, r *http.Request) {

	serv.Page.Title = "Добавление блога"
	serv.Page.Data = models.Blog{}
	serv.Page.Command = "new"

	if err := serv.dictionary["BLOG"].ExecuteTemplate(w, "base", serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// addBlogHandler - добавляет блог
// @Summary Добавляет блог в список
// @Description добавить блог в список
// @Tags blog, add
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /api/v1/blog [post]
func (serv *Server) addBlogHandler(w http.ResponseWriter, r *http.Request) {

	data, _ := ioutil.ReadAll(r.Body)

	blog := models.Blog{}
	_ = json.Unmarshal(data, &blog)

	if err := blog.AddBlog(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	http.Redirect(w, r, "/", 301)
}

// deleteBlogHandler - удаляет блог
// @Summary Удаляет блог из список
// @Description удаляет блог из список
// @Tags blog, delete
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /api/v1/blog/{id} [delete]
func (serv *Server) deleteBlogHandler(w http.ResponseWriter, r *http.Request) {

	blogID, err := getBlogIDFromRequest(r)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	blog := models.Blog{}
	blog.ID = blogID

	if err := blog.Delete(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	http.Redirect(w, r, "/", 301)
}
