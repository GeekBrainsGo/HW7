package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"serv/server"

	"github.com/sirupsen/logrus"
)

// @title GeekBrains HW7 Server
// @version 1.0
// @description Server for api & documentation testing
// @contact.name Yuri Kulagin
// @contact.url https://t.me/jkulvich
// @contact.email jkulvichi@gmail.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1

func main() {
	flagAddr := flag.String("addr", "localhost:8080", "server address")
	flagSwagURLAddr := flag.String("swag-url", "http://localhost:8080/api/v1/docs/swagger.json", "swagger.json http address")
	flagSwagJSONPath := flag.String("swag-path", "docs/swagger.json", "swagger.json path")
	flag.Parse()

	ctx := context.Background()
	lg := logrus.New()

	serv := server.NewServer(ctx, lg)
	serv.SetSwagger(*flagSwagJSONPath, *flagSwagURLAddr)
	serv.Start(*flagAddr)

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan

	serv.Stop()
}
