package mongo

type Config struct {
	URI         string
	DBName      string
	EnableTrace bool
}

var Conf map[string]*Config
