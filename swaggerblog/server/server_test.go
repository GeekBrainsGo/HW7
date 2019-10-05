package server_test

import (
	"HW7/swaggerblog/models"
	"HW7/swaggerblog/server"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testCases = []struct {
	description string
	method      string
	post        models.Post
}{
	{
		"post test",
		http.MethodPost,
		models.Post{Title: "1", Author: "2", Content: "3"},
	},
}

func TestPostHandle(t *testing.T) {
	msg := `
	Description: %s
	Expected: %q
	Got: %q
	`
	lg := logrus.New()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		lg.WithError(err).Fatal("can't get new client")
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}

	if err = client.Ping(ctx, nil); err != nil {
		lg.WithError(err).Fatal("can't ping to db")
	}

	db := client.Database("mongoblog")

	testserver := server.New(ctx, lg, "./www", db)
	for _, test := range testCases {
		j, _ := json.Marshal(test.post)

		request, _ := http.NewRequest(test.method, "/api/v1/", bytes.NewReader(j))
		response := httptest.NewRecorder()
		testserver.PostHandle(response, request)

		post := &models.Post{}
		json.NewDecoder(response.Body).Decode(post)
		request.Body.Close()
		response.Body.Reset()

		if post.Title != test.post.Title {
			t.Fatalf(msg, test.description, test.post.Title, post.Title)
		}
		if post.Author != test.post.Author {
			t.Fatalf(msg, test.description, test.post.Author, post.Author)
		}
		if post.Content != test.post.Content {
			t.Fatalf(msg, test.description, test.post.Content, post.Content)
		}
	}
	err = testserver.Stop() // TODO check why panics
	client.Disconnect(ctx)
}

func TestDeleteHandle(t *testing.T) {
	// TODO implement deletion check
}

func TestPutHandle(t *testing.T) {
	// TODO implement put check
}
