package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fedorwk/ngrok-single-endpoint/config"
	"github.com/fedorwk/ngrok-single-endpoint/server"
	"github.com/fedorwk/ngrok-single-endpoint/srvlist"
)

func main() {
	config.Init()
	srvlistFile, err := os.Open(config.Cfg.ServiceListFilePath)
	if err != nil {
		log.Fatalf("err opening srvlist file: %+v\n", err)
	}
	srvlist, err := srvlist.FromReader(srvlistFile)
	if err != nil {
		log.Fatalf("err reading services: %+v\n", err)
	}
	client := http.DefaultClient
	httpserver := http.Server{
		Addr: ":" + config.Cfg.Port,
	}

	server := server.New(&httpserver, client, srvlist)
	err = server.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
