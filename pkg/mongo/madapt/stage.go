package madapt

import (
	"strings"

	"github.com/wolif/gosaber/pkg/ref"
	"github.com/wolif/gosaber/pkg/ref/structtags"
	"github.com/wolif/gosaber/pkg/util/strs"
	"go.mongodb.org/mongo-driver/bson"
)

type matchResult bson.M

const SearchTag = "search"

func Match(cond interface{}) bson.M {
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
		return bson.M{}
	}
	m := make(matchResult)
	condRef := ref.New(cond)
	for fieldName, options := range tags {
		val, _ := condRef.StructValueGet(fieldName)
		searchField := strings.TrimSpace(options[fieldToken])
		if searchField == "" {
			searchField = strs.SnakeCase(fieldName)
		}
		m.matchOpField(searchField, options["op"], val, options["omitEmpty"], options["addOp"])
	}

	return bson.M(m)
}

func Sort(fields ...string) bson.D {
	ret := bson.D{}
	if len(fields) == 0 {
		return ret
	}

	for _, field := range fields {
		if strings.HasPrefix(field, "-") {
			ret = append(ret, bson.E{Key: field[1:], Value: -1})
		} else {
			ret = append(ret, bson.E{Key: field, Value: 1})
		}
	}
	return ret
}
