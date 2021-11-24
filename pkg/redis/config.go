package redis

import (
	"fmt"
	goredis "github.com/go-redis/redis"
)

type Config struct {
	Addrs []string `json:"addrs"`
	Pwd   string   `json:"pwd"`
	DB    int      `json:"db"`
}

var Client goredis.Cmdable
var Nil = goredis.Nil

func Init(config *Config) error {
	switch {
	case len(config.Addrs) == 1:
		Client = goredis.NewClient(
			&goredis.Options{
				Addr:     config.Addrs[0], // use default Addr
				Password: config.Pwd,      // no password set
				DB:       config.DB,       // use default DB
			})
	case len(config.Addrs) > 1:
		Client = goredis.NewClusterClient(
			&goredis.ClusterOptions{
				Addrs:    config.Addrs,
				Password: config.Pwd,
			})
	default:
		return fmt.Errorf("redis addrs is empty")
	}

	if err := Client.Ping().Err(); err != nil {
		return err
	}

	return nil
}
