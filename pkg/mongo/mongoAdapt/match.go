package mongoAdapt

import (
	"go.mongodb.org/mongo-driver/bson"
	"github.com/wolif/gosaber/pkg/ref"
)

const (
	TagSrhField  = "search"
	TagCmpOpt    = "opt"
	TagOmitEmpty = "omitEmpty"
	TagAddOpt    = "addOpt"
)

type matchResult bson.M

func Match(cond interface{}) bson.M {
	condRef := ref.New(cond)
	if !condRef.IsStruct() {
		return bson.M{}
	}

	m := make(matchResult)
	fields, _ := condRef.GetStructFields()
	for fieldName, _ := range fields {
		searchField, _ := condRef.GetStructFieldTag(fieldName, TagSrhField)
		opt, _ := condRef.GetStructFieldTag(fieldName, TagCmpOpt)
		omitempty, _ := condRef.GetStructFieldTag(fieldName, TagOmitEmpty)
		addopt, _ := condRef.GetStructFieldTag(fieldName, TagAddOpt)

		if searchField == "" {
			continue
		}

		fieldValue, _ := condRef.GetStructFieldValue(fieldName)

		if opt == "" {
			m.matchField(searchField, fieldValue, omitempty, addopt)
		} else {
			m.matchOptField(searchField, opt, fieldValue, omitempty, addopt)
		}
	}

	return (bson.M)(m)
}
