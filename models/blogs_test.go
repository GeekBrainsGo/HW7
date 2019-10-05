package models

import (
	"testing"

	_ "github.com/go-sql-driver/MySQL"
	"github.com/jinzhu/gorm"
)

const (
	strSQL   = "root:root@/MyBlogs?parseTime=true"
	strTitle = "Это мой первый блог"
)

func connectDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open("mysql", strSQL)
	if err != nil {
		t.Errorf("can't connect to db, ERROR: %v", err)
	}
	return db
}

func TestAddBlog(t *testing.T) {

	db := connectDB(t)
	defer db.Close()

}

func TestGetBlog(t *testing.T) {

	db := connectDB(t)
	defer db.Close()

	blog, err := GetBlog(db, 1)
	if err != nil {
		t.Errorf("can't get record from db, ERROR: %v", err)
	}

	if blog.Title != strTitle {
		t.Errorf("unexpected Title blog, got: %v, want %v", blog.Title, strTitle)
	}
}

func TestDeleteBlog(t *testing.T) {

	db := connectDB(t)
	defer db.Close()

	blog, err := GetBlog(db, 1)
	if err != nil {
		t.Errorf("can't get record from db, ERROR: %v", err)
	}

	db.Delete(&blog)

	blog, err = GetBlog(db, 1)
	if err == nil {
		t.Errorf("blob has not been deleted")
	}

}
