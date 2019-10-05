package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

)

// Post - объект поста в блоге
type Post struct {
	Mongo   `inline`
	//ID        uint `bson:"id" json:"id"`
	Title     string `bson:"title" json:"title"`
	Text      string `bson:"text" json:"text"`
	Labels    []string `bson:"labels" json:"labels"`
}

// GetMongoCollectionName - Перегруженный метод возвращающий имя коллекции структуры
func (p *Post) GetMongoCollectionName() string {
	return "posts"
}

// Posts - массив постов в блоге
type Posts []Post

// Create - создает задачу в БД
func (p *Post) Create(ctx context.Context, db *mongo.Database) (*Post, error) {
	col := db.Collection(p.GetMongoCollectionName())
	res, err := col.InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}
	p.ID = res.InsertedID.(primitive.ObjectID)
	return p, nil
}

// Delete - удалить объект из базы
func (p *Post) Delete(ctx context.Context, db *mongo.Database) (*Post, error) {
	col := db.Collection(p.GetMongoCollectionName())
	_, err := col.DeleteOne(ctx, bson.M{"_id": p.ID})
	return p, err
}

// Update - обновляет объект в БД
func (p *Post) Update(ctx context.Context, db *mongo.Database) (*Post, error) {
	col := db.Collection(p.GetMongoCollectionName())
	res, err := col.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, errors.New("Not found post")
	}

	return p, nil
}

// GetPost - получение поста
func GetPost(postId primitive.ObjectID, ctx context.Context, db *mongo.Database) (*Post, error) {
	p := &Post{
		Mongo:  Mongo{
			ID:   postId,
		},
	}

	col := db.Collection(p.GetMongoCollectionName())
	err := col.FindOne(ctx, bson.M{"_id": p.ID}).Decode(&p)

	return p, err
}

// GetPosts - Возвращает все посты
func GetPosts(ctx context.Context, db *mongo.Database) ([]Post, error) {
	p := Post{}
	col := db.Collection(p.GetMongoCollectionName())

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// Find - Находит все посты у которых поле равно значению
func Find(ctx context.Context, db *mongo.Database, field string, value interface{}) ([]Post, error) {
	p := Post{}
	col := db.Collection(p.GetMongoCollectionName())

	cur, err := col.Find(ctx, bson.M{field: value})
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}