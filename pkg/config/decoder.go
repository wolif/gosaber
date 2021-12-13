package config

import (
	"flag"
	"github.com/wolif/gosaber/pkg/config/envvariant"
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
