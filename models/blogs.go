package models

import (
	"github.com/jinzhu/gorm"
)

// Blog - объект блога
type Blog struct {
	gorm.Model
	Title string `gorm:"column:title;type:varchar(45)"`
	Body  string `json:"article" gorm:"column:body;default:''"`
}

// // BlogItems - список блогов
// type BlogItems []Blog

// TableName - имя таблицы блогов
func (Blog) TableName() string {
	return "blogs"
}

// AddBlog - обновляет объект в БД
func (blog *Blog) AddBlog(db *gorm.DB) error {

	db.Create(&blog)

	return nil
}

// UpdateBlog - обновляет объект в БД
func (blog *Blog) UpdateBlog(db *gorm.DB) error {

	db.Save(&blog)

	return nil
}

// GetAllBlogs - получение всех блогов
func GetAllBlogs(db *gorm.DB) ([]Blog, error) {

	blogs := []Blog{}
	db.Find(&blogs)

	return blogs, nil
}

// GetBlog - получение блогa
func GetBlog(db *gorm.DB, id uint) (Blog, error) {

	blog := Blog{}

	// Get record with primary key (only works for integer primary key)
	db.First(&blog, id)

	return blog, nil
}

// Delete - удалят объект из базы
func (blog *Blog) Delete(db *gorm.DB) error {

	// soft deleted
	db.Delete(&blog) // UPDATE blogs SET deleted_at="2013-10-29 10:23" WHERE id = NUM;
	return nil
}
