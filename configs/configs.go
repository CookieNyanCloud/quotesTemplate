package configs

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

const (
	tgToken = "TG_TOKEN"
	url     = "URL"
)

type Conf struct {
	TgToken string
	URL     string
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
		URL:     os.Getenv(url),
	}, nil
}
