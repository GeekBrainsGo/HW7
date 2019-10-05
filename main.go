package main

import (
	"flag"
	"os"
	"os/signal"
	"serv/server"

	_ "github.com/go-sql-driver/MySQL"
	"github.com/jinzhu/gorm"
)

// @title Мой первый простой блог
// @version 2.0
// @description Этот блог создан для учебных целей. Использует MySQL для хранения блогов.

// @contact.email vbaykin@mail.ru

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	flagRootDir := flag.String("rootdir", "./web", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	flagMySQL := flag.String("sql", "root:root@/MyBlogs?parseTime=true", "MySQL connection string, format: user:password@tcp(host:port)/database")
	flag.Parse()

	lg := NewLogger()

	db, err := gorm.Open("mysql", *flagMySQL)
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err.Error()) // proper error handling instead of panic in your app
	// }

	serv := server.New(lg, *flagRootDir, db)

	go func() {
		err := serv.Start(*flagServAddr)
		if err != nil {
			lg.WithError(err).Fatal("can't run the server")
		}
	}()

	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt, os.Kill)
	<-stopSig
	lg.Info("Stop server!")

}
