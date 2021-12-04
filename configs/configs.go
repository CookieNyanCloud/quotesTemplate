package configs

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
)

const (
	tgToken = "TG_TOKEN"
	admin   = "ADMIN"
)

type Conf struct {
	TgToken string
	Admin   string
}

func InitConf() (*Conf, error) {
	var local bool
	flag.BoolVar(&local, "local", false, "хост")
	flag.Parse()
	return envVar(local)
}

func envVar(local bool) (*Conf, error) {
	if local {
		err := godotenv.Load(".env")
		if err != nil {
			return &Conf{}, err
		}
	}
	return &Conf{
		TgToken: os.Getenv(tgToken),
		Admin:   os.Getenv(admin),
	}, nil
}
