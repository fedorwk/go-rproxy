package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	Cfg     config
	cfgPath string
)

// TODO func config.Get() *Cfg -> cfg (err if not initialized)

type config struct {
	// Server section
	Port string `yaml:"port" env:"PORT" env-default:"8080"`

	// App Section
	ServiceListFilePath string `yaml:"srvlistpath" env:"SRVLISTPATH"`
}

func Init() error {
	flagSet := flag.NewFlagSet("Main", flag.ContinueOnError)
	flagSet.StringVar(&cfgPath, "cfg", "config.yml", "path to config file")
	flagSet.Usage = cleanenv.FUsage(flagSet.Output(), &Cfg, nil, flagSet.Usage)

	flagSet.Parse(os.Args[1:])
	err := cleanenv.ReadConfig(cfgPath, &Cfg)
	if err != nil {
		return err
	}
	return nil
}
