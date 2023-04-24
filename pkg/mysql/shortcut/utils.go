package shortcut

import (
	"github.com/jinzhu/gorm"
	"github.com/wolif/gosaber/pkg/mysql"
)

var connName = "default"

func SetConnName(cn string) {
	connName = cn
}

func resolveDB(db ...*gorm.DB) (*gorm.DB, error) {
	if len(db) > 0 && db[0] != nil {
		return db[0], nil
	}
	return mysql.GetDbByName(connName)
}

// -----------------------------------------------------------------------------
const (
	DefPage     int = 1
	DefPageSize int = 10
)

func resolvePagi(page, pageSize int, defOption ...int) (int, int) {
	defPage := DefPage
	defPageSize := DefPageSize
	if len(defOption) >= 1 && defOption[0] > 0 {
		defPage = defOption[0]
	}
	if len(defOption) >= 2 && defOption[1] > 0 {
		defPageSize = defOption[1]
	}

	if page <= 0 {
		page = defPage
	}
	if pageSize <= 0 {
		pageSize = defPageSize
	}
	return page, pageSize
}