package helper

import (
	"github.com/jinzhu/gorm"
	"github.com/wolif/gosaber/pkg/mysql/where"
)

func WithWhere(db *gorm.DB, whereCond interface{}) *gorm.DB {
	wv := where.Make(whereCond)
	if wv.HasCondition() {
		db = db.Where(wv.Sql, wv.Value...)
	}
	return db
}
