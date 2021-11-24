package gormhelper

import "github.com/jinzhu/gorm"

const (
	DefPage     int = 1
	DefPageSize int = 10
)

func Paginate(db *gorm.DB, page, pageSize int, defOption ...int) *gorm.DB {
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

	return db.Offset((page - 1) * pageSize).Limit(pageSize)
}
