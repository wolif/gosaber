package shortcut

import (
	"github.com/jinzhu/gorm"
	"github.com/wolif/gosaber/pkg/mysql"
)

type transaction struct {
	TX *gorm.DB
}

func NewTransaction() (*transaction, error) {
	db, err := mysql.GetDbByName(connName)
	if err != nil {
		return nil, err
	}
	return &transaction{TX: db}, nil
}

func (t *transaction) Begin() *gorm.DB {
	return t.TX.Begin()
}

func (t *transaction) Rollback() {
	t.TX.Rollback()
	t = nil
}

func (t *transaction) Commit() error {
	err := t.TX.Commit().Error
	t = nil
	return err
}

func (t *transaction) Do(fn func(*gorm.DB) error) error {
	tx := t.Begin()
	defer func() {
		if e := recover(); e != nil {
			t.Rollback()
		}
	}()
	err := fn(tx)
	if err != nil {
		t.Rollback()
		return err
	}
	return t.Commit()
}
