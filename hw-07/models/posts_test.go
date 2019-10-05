package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

var CTX = context.Background()
var db = getDB("mongodb://localhost:27017")

func getDB(uri string) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	_ = client.Connect(CTX)
	db := client.Database("geekbrains")
	return db
}

func TestPostItem_Insert(t *testing.T) {
	post := PostItem{
		Mongo:            Mongo{},
		Idn:              "5",
		Title:            "5",
		Date:             "5",
		SmallDescription: "5",
		Description:      "5",
	}
	_, err := post.Insert(CTX, db)
	if err != nil {
		t.Error("insertion of the record into the database failed")
	}
}

func TestPostItem_Delete(t *testing.T) {
	post := PostItem{
		Mongo:            Mongo{},
		Idn:              "5",
		Title:            "5",
		Date:             "5 ",
		SmallDescription: "5",
		Description:      "5",
	}
	_, err := post.Delete(CTX, db)
	if err != nil {
		t.Error("database record deletion failed")
	}
}

func TestPostItem_Update(t *testing.T) {
	post := PostItem{
		Mongo:            Mongo{},
		Idn:              "4",
		Title:            "44",
		Date:             "44",
		SmallDescription: "44",
		Description:      "44",
	}
	_, err := post.Update(CTX, db)
	if err != nil {
		t.Error("database record update failed")
	}
}
