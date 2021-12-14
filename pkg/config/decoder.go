package config

import (
	"flag"
	"github.com/wolif/gosaber/pkg/config/envvariant"
	"github.com/wolif/gosaber/pkg/util/path"
)

var (
	Path     string
)

func init() {
	flag.StringVar(&Path, "conf", "", "config path")
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
