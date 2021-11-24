package gormhelper

import (
	"github.com/jinzhu/gorm"
)

func scopeNullFunc(db *gorm.DB) *gorm.DB {
	return db
}
