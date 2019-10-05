package main

/*
	Basics Go.
	Rishat Ishbulatov, dated Oct 03, 2019.
	Write tests for the functions with which you work with the database.
	Write http tests for site methods.
	Add documentation to the project.
	Add a Makefile to Easily Deploy Your Server.
*/

import (
	"context"
	"flag"
	"go_basics/packages/swaggerblog/server"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Mongoblog server
// @version 1.0
// @description Blog server with swagger and tests for handlers.
// @contact.name Rishat Ishbulatov
// @contact.email progjb@gmail.com
// @license.name MIT
// @host localhost:8080
// @BasePath /

func main() {
	flagRootDir := flag.String("root", "./www", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	flagDBname := flag.String("dbname", "mongoblog", "database name")
	flagSwagURLAddr := flag.String("swag-url", "http://localhost:8080/api/v1/docs/swagger.json", "swagger.json http address")
	flagSwagJSONPath := flag.String("swag-path", "docs/swagger.json", "swagger.json path")
	flag.Parse()

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

	db := client.Database(*flagDBname)
	defer client.Disconnect(ctx)

	serv := server.New(ctx, lg, *flagRootDir, db)
	serv.SetSwagger(*flagSwagJSONPath, *flagSwagURLAddr)
	serv.Start(*flagServAddr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	serv.Stop()
}
