package mysql

type Config struct {
	DbUrl           string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	LogMode         bool
}
