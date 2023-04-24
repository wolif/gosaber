package where

import (
	"strings"

	"github.com/wolif/gosaber/pkg/ref"
	"github.com/wolif/gosaber/pkg/ref/structtags"
	"github.com/wolif/gosaber/pkg/util/strs"
)

const SearchTag = "search"

type WhereVal struct {
	Sql   string
	Value []interface{}
}

func (wv *WhereVal) HasCondition() bool {
	return len(wv.Value) > 0
}

func Make(cond interface{}) *WhereVal {
	ret := &WhereVal{
		Sql:   "",
		Value: make([]interface{}, 0),
	}
	fieldToken := ""
	tags, err := structtags.Parse(
		cond,
		SearchTag,
		map[string]string{
			fieldToken:  "",
			"op":        "",
			"omitEmpty": "true",
			"addOp":     "",
		},
	)
	if err != nil {
		return ret
	}

	condRef := ref.New(cond)
	for fieldName, options := range tags {
		val, _ := condRef.StructValueGet(fieldName)
		searchField := strings.TrimSpace(options[fieldToken])
		if searchField == "" {
			searchField = strs.SnakeCase(fieldName)
		}
		where(ret, searchField, options["op"], val, options["omitEmpty"], options["addOp"])
	}
	if ret.HasCondition() {
		ret.Sql = strings.TrimLeft(ret.Sql, " AND")
	}
	return ret
}
