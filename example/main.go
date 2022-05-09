package main

import (
	"log"
	"net/http"
	"os"

	rproxy "github.com/fedorwk/go-rproxy"
	"github.com/gin-gonic/gin"
)

func main() {
	srvlistFile, err := os.Open("srvlist")
	if err != nil {
		log.Fatalf("err opening srvlist file: %+v\n", err)
	}
	srvlist, err := rproxy.NewServiceListFromReader(srvlistFile)
	if err != nil {
		log.Fatalf("err reading services: %+v\n", err)
	}

	proxyRouter, err := rproxy.New(gin.New(), srvlist)
	if err != nil {
		log.Fatalln("setting proxy router", err)
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: proxyRouter,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
