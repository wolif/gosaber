package helper

import "github.com/jinzhu/gorm"

const (
	DefPage     int = 1
	DefPageSize int = 10
)

func WithPaginate(db *gorm.DB, page, pageSize int, defOption ...int) *gorm.DB {
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

func PgInDef(db *gorm.DB, page, pageSize int, defOption ...int) *gorm.DB {
	return WithPaginate(db, page, pageSize, defOption...)
}

func PgOutDef(db *gorm.DB, pg ...int) *gorm.DB {
	if len(pg) == 2 && pg[0] > 0 && pg[1] > 0 {
		return db.Offset((pg[0] - 1) * pg[1]).Limit(pg[1])
	}
	return db
}
