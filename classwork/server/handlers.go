package server

import (
	"io/ioutil"
	"net/http"
)

// HandleSwagger - Returns swagger.json docs
// @Description Returns swagger.json docs
// @Tags system
// @Success 200 {string} string
// @Router /docs/swagger.json [get]
func (serv *Server) HandleSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, serv.swagPath)
}

// HandleEcho - Returns echo
// @Description Returns echo
// @Tags health-check
// @Param data body string true "Any text"
// @Success 200 {string} string
// @Failure 500 {object} models.ServErr
// @Router /echo [post]
func (serv *Server) HandleEcho(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	w.Write(data)
}
