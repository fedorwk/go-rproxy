package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	rproxy "github.com/fedorwk/ngrok-single-endpoint"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	Config.Init()
	srvlistFile, err := os.Open(Config.ServiceListFilePath)
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
		Addr:    ":" + Config.Port,
		Handler: proxyRouter,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

var (
	Config  config
	cfgPath string
)

type config struct {
	// Server section
	Port string `yaml:"port" env:"PORT" env-default:"8080"`

	// App Section
	ServiceListFilePath string `yaml:"srvlistpath" env:"SRVLISTPATH"`
}

func (c *config) Init() error {
	flagSet := flag.NewFlagSet("Main", flag.ContinueOnError)
	flagSet.StringVar(&cfgPath, "cfg", "config.yml", "path to config file")
	flagSet.Usage = cleanenv.FUsage(flagSet.Output(), c, nil, flagSet.Usage)

	flagSet.Parse(os.Args[1:])
	err := cleanenv.ReadConfig(cfgPath, c)
	if err != nil {
		return err
	}
	return nil
}
