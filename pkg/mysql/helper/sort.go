package helper

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

/*
 * usage
 * "id"  -> " ORDER BY `id` ASC "
 * "-id" -> " ORDER BY `id` DESC "
 */
func WithSort(db *gorm.DB, sort ...string) *gorm.DB {
	for _, s := range sort {
		if s = strings.TrimSpace(s); s != "" {
			db = db.Order(makeSortString(s))
		}
	}
	return db
}

func makeSortString(s string) string {
	ret := ""
	if sbs := []byte(s); sbs[0] == '-' {
		ret = fmt.Sprintf("`%s` DESC", string(sbs[1:]))
	} else {
		ret = fmt.Sprintf("`%s` ASC", s)
	}
	return ret
}
