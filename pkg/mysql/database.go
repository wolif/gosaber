package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

var (
	connPool = make(map[string]*gorm.DB)
	nutex    sync.Mutex
	DefConn  = "default"
)

func GetDefaultDb() (*gorm.DB, error) {
	return GetDbByName(DefConn)
}

func InitDb(name string, conf *Config) error {
	nutex.Lock()
	defer nutex.Unlock()

	maxIdle := 5
	if conf.MaxIdleConns > 0 {
		maxIdle = conf.MaxIdleConns
	}

	maxOpenConns := 10
	if conf.MaxOpenConns > 0 {
		maxOpenConns = conf.MaxOpenConns
	}

	connMaxLifetime := time.Duration(3600)
	if conf.ConnMaxLifetime > 0 {
		connMaxLifetime = time.Duration(conf.ConnMaxLifetime)
	}

	db, err := gorm.Open("mysql", conf.DbUrl)
	if err != nil {
		return err
	}

	db.LogMode(conf.LogMode)
	db.DB().SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpenConns)

	if err = db.DB().Ping(); err != nil {
		return err
	}

	connPool[name] = db

	return nil
}

func GetDbByName(name string) (*gorm.DB, error) {
	nutex.Lock()
	defer nutex.Unlock()

	if connPool == nil {
		return nil, fmt.Errorf("db pool is empty")
	}

	if v, ok := connPool[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("%s not found in db pool", name)
}
