package configs

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

const (
	tgToken = "TG_TOKEN"
	addr    = "ADDR"
)

type Conf struct {
	TgToken string
	Addr    string
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
		Addr:    os.Getenv(addr),
	}, nil
}
