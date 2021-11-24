package gormhelper

import (
	"github.com/jinzhu/gorm"
	"github.com/wolif/gosaber/pkg/ref"
)

const (
	SearchFieldTag       string = "search"
	OmitEmptyValueTag    string = "omitEmpty"
	CompareOperationTag  string = "opt"
	AdditionOperationTag string = "addOpt"
)

func Scope(db *gorm.DB, searchCond interface{}) *gorm.DB {
	condRef := ref.New(searchCond)
	if condRef.IsStruct() {
		fields, _ := condRef.GetStructFields()
		for fieldName, _ := range fields {
			searchField, _ := condRef.GetStructFieldTag(fieldName, SearchFieldTag)
			omitEmpty, _ := condRef.GetStructFieldTag(fieldName, OmitEmptyValueTag)
			opt, _ := condRef.GetStructFieldTag(fieldName, CompareOperationTag)
			additionOpt, _ := condRef.GetStructFieldTag(fieldName, AdditionOperationTag)

			if searchField == "" {
				continue
			}

			fieldValue, _ := condRef.GetStructFieldValue(fieldName)
			if opt == "" {
				db = db.Scopes(where(searchField, fieldValue, omitEmpty, additionOpt))
			} else {
				db = db.Scopes(whereOpt(searchField, opt, fieldValue, omitEmpty, additionOpt))
			}
		}
	}

	return db
}
