package main

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleEditPostA(t *testing.T) {

	lg := logrus.New()
	ctx := context.Background()
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	_ = client.Connect(ctx)
	db := client.Database("blog")

	serv := NewServer(lg, db, ctx)
	serv.Start(":8080")

	var jsonStr = []byte(`{"ID":"5d9886508b09c75bc0ec8b6e","title":"Пост 1","text":"Новый очень интересный текст","labels":["l1","l2"]}`)
	req, err := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serv.HandleEditPost)

	handler.ServeHTTP(rr, req)

	// Проверяем код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := `{"ID":"5d9886508b09c75bc0ec8b6e","title":"Пост 1","text":"Новый очень интересный текст","labels":["l1","l2"]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serv.Stop()
}

func TestHandleEditPostB(t *testing.T) {

	lg := logrus.New()
	ctx := context.Background()
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	_ = client.Connect(ctx)
	db := client.Database("blog")

	serv := NewServer(lg, db, ctx)
	serv.Start(":8080")

	var jsonStr = []byte(`{"ID":"5d9886508b09c75bc0ec8b6e","title":"Пост 1","text":"Изменен","labels":["l1"]}`)
	req, err := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serv.HandleEditPost)

	handler.ServeHTTP(rr, req)

	// Проверяем код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := `{"ID":"5d9886508b09c75bc0ec8b6e","title":"Пост 1","text":"Изменен","labels":["l1"]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serv.Stop()
}
