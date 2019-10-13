package models

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetBlogs(t *testing.T) {
	db := CreateMongoDb()
	_, err := GetBlogs(nil, db)
	if err != nil {
		t.Error(err)
	}
}

func TestBlog_Insert(t *testing.T) {
	db := CreateMongoDb()
	testCases := Blogs{
		Blog{
			Title:    "test1",
			Contents: "test1"},
		Blog{
			Title:    "test2",
			Contents: "test2"},
		Blog{
			Title:    "test3",
			Contents: "test3"},
	}
	for _, tcase := range testCases {
		_, err := tcase.Insert(nil, db)
		if err != nil {
			t.Error("can't insert Blog item into DB: ", err)
		}
	}
}

func TestFindBlog(t *testing.T) {
	db := CreateMongoDb()
	var testCases []string
	blogs, _ := GetBlogs(nil, db)
	for _, blog := range blogs {
		testCases = append(testCases, blog.ID.Hex())
	}
	for _, tcase := range testCases {
		blog, err := FindBlog(nil, db, tcase)
		if err != nil {
			t.Error("Can't get object from DB: ", err)
		} else if blog.ID.Hex() != tcase {
			t.Errorf("Expected blog id: %s, but got id: %s", tcase, blog.ID)
		}
	}
}

func TestBlog_Update(t *testing.T) {
	db := CreateMongoDb()
	tcases, _ := GetBlogs(nil, db)
	for _, tcase := range tcases {
		restCase, err := tcase.Update(nil, db)
		if tcase.ID != restCase.ID {
			t.Errorf("ID's mismatch, expected: %s, got %s:", tcase, restCase)
		}
		if err != nil {
			t.Error("error on updating: ", err)
		}
	}
}

func TestBlog_Delete(t *testing.T) {
	db := CreateMongoDb()
	tcases, _ := GetBlogs(nil, db)
	for _, tcase := range tcases {
		_, err := tcase.Delete(nil, db)
		if err != nil {
			t.Error("error on updating: ", err)
		}
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
