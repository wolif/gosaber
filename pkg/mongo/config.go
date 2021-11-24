package mongo

type Config struct {
	URI     string
	DBName  string
}

var Conf map[string]*Config
