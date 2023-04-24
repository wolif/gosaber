package config

import (
	"flag"

	"github.com/wolif/gosaber/pkg/config/dotenv"
	"github.com/wolif/gosaber/pkg/config/json"
	"github.com/wolif/gosaber/pkg/config/properties"
	"github.com/wolif/gosaber/pkg/util/path"
)

var (
	Path string
)

func init() {
	flag.StringVar(&Path, "conf", "", "config path")
}

func LoadProperties(conf interface{}) error {
	var envPath string
	if Path == "" {
		p, err := path.Find("configs/.properties", 5)
		if err != nil {
			return err
		}
		envPath = p
	} else {
		envPath = Path
	}

	return properties.Load(envPath, conf)
}

func LoadJson(conf interface{}) error {
	var jsonPath string
	if Path == "" {
		p, err := path.Find("configs/config.json", 5)
		if err != nil {
			return err
		}
		jsonPath = p
	} else {
		jsonPath = Path
	}
	return json.Load(jsonPath, conf)
}

func LoadEnv(conf interface{}) error {
	var dotenvPath string
	if Path == "" {
		p, err := path.Find("configs/.env", 5)
		if err != nil {
			return err
		}
		dotenvPath = p
	} else {
		dotenvPath = Path
	}
	return dotenv.Load(dotenvPath)
}
