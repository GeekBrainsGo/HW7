package server

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"serv/logger"
	"serv/models"
	"testing"
)

func TestServer_postNewHandler(t *testing.T) {
	flagRootDir := flag.String("rootdir", "./www", "root dir of the server")
	flag.Parse()

	ctx := context.Background()
	lg := logger.NewLogger()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	_ = client.Connect(ctx)
	db := client.Database("geekbrains")
	serv := New(lg, *flagRootDir, ctx, db)

	post := models.PostItem{
		Title:            "title_1",
		Date:             "2019-10-04",
		SmallDescription: "SmallDescription_1",
		Description:      "Description_1",
	}

	data, _ := json.Marshal(post)
	reader := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", "/api/v1/posts/new", reader)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serv.postNewHandler)
	handler.ServeHTTP(rr, req)

	resp := rr.Body
	data, _ = ioutil.ReadAll(resp)
	postNew := models.PostItem{}
	_ = json.Unmarshal(data, &postNew)

	if postNew.Title != post.Title ||
		postNew.SmallDescription != post.SmallDescription ||
		postNew.Description != post.Description ||
		postNew.Date != post.Date {
		t.Errorf("add new post failed")
	}

	fmt.Println(rr.Body)

}
