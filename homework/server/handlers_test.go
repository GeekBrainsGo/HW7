package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sirupsen/logrus"
)

func TestHandleGetIndex(t *testing.T) {
	lg := logrus.New()
	lg.SetLevel(0)
	db := CreateMongoDb()
	serv := New(lg, db, "www/static")
	reader := bytes.NewReader([]byte("TEST"))
	req, _ := http.NewRequest("GET", "/", reader)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serv.handleGetIndex)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error("Status from server: ", req.Response.Status)
	}

}

func CreateMongoDb() *mongo.Database {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)
	ctx := context.Background()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		lg.Panic("Can't connect to DB", err)
		return nil
	}
	_ = client.Connect(ctx)
	db := client.Database("blogs")
	return db
}
