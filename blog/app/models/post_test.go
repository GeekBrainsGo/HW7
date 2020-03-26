package models

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const POSTID = "testPostID"

func TestPost_Create(t *testing.T) { // 3

	var newPost Post

	newPost.ID = POSTID
	newPost.Title = "New Post Title"
	newPost.Short = "Short body"
	newPost.Body = "Content body"

	session, err := mgo.Dial("mongodb://root:123qweqwe@127.0.0.1:27017")
	if err != nil {
		t.Fatal(err)
	}

	conn := session.DB("blog").C("posts")

	err = conn.Insert(newPost)
	if err != nil {
		t.Fatal(err)
	}

}

func TestPost_Find(t *testing.T) {

	session, err := mgo.Dial("mongodb://root:123qweqwe@127.0.0.1:27017")
	if err != nil {
		t.Fatal(err)
	}

	conn := session.DB("blog").C("posts")

	var post Post

	err = conn.Find(bson.M{"id": POSTID}).One(&post)
	if err != nil {
		t.Fatal(err)
	}

}

func TestPost_Delete(t *testing.T) { // 3

	session, err := mgo.Dial("mongodb://root:123qweqwe@127.0.0.1:27017")
	if err != nil {
		t.Fatal(err)
	}

	conn := session.DB("blog").C("posts")

	err = conn.Remove(bson.M{"id": POSTID})
	if err != nil {
		t.Fatal(err)
	}

}
