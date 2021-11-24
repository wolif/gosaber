package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/wolif/gosaber/pkg/config/envvariant"
	"github.com/wolif/gosaber/pkg/util/env"
	"github.com/wolif/gosaber/pkg/util/path"
)

var (
	Path     string
	confName string
	Test     = flag.Bool("t", false, "test conf file")
)

func init() {
	flag.StringVar(&Path, "conf", "", "config path")
}

func LoadToml(conf interface{}) error {
	var tomlPath string
	if Path == "" {
		confName = env.GetEnvWithFallback("CONF_NAME", "example")
		p, err := path.FindPath(fmt.Sprintf("configs/%s.toml", confName), 5)
		if err != nil {
			return err
		}
		tomlPath = p
	} else {
		tomlPath = Path
	}

	_, err := toml.DecodeFile(tomlPath, conf)
	return err
}

func LoadEnv(conf interface{}) error {
	var envPath string
	if Path == "" {
		p, err := path.FindPath("configs/.env", 5)
		if err != nil {
			return err
		}
		envPath = p
	} else {
		envPath = Path
	}

	return envvariant.Load(envPath, conf)
}
