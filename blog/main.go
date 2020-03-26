package main

import (
	"blog/app/webserver"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() { // 1
	flag.StringVar(&configPath, "config-path", "configs/blog.toml", "path to config file")
}

func main() { // 3
	flag.Parse()

	config := webserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := webserver.Start(config); err != nil {
		log.Fatal(err)
	}

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan

}
