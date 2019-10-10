package main

// @title HW7
// @version 0.1
// @description This is a HW blog

// @contact.name Dmitrii Fadeev

import (
	"context"
	"flag"
	"serv/models"
	"serv/server"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const staticDir = "www/static"

func main() {
	flagServAddr := flag.String("addr", ":8080", "server address")
	flagConnDb := flag.String("conndb", "mongodb://localhost:27017", "db conn string")

	lg := NewLogger()
	ctx := context.Background()
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagConnDb))
	if err != nil {
		lg.Panic("Can't connect to DB", err)
	} else {
		lg.Info("Connection to DB successful")
	}

	_ = client.Connect(ctx)
	db := client.Database("blogs")

	blog := &models.Blog{
		Title:    "test",
		Contents: "test",
	}

	_, err = blog.Insert(ctx, db)
	if err != nil {
		lg.Fatal(err)
	}

	srv := server.New(lg, db, staticDir)
	srv.Start(*flagServAddr)
}
