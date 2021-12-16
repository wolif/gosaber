package structtags

import (
	"strings"

	"github.com/wolif/gosaber/pkg/ref"
	"github.com/wolif/gosaber/pkg/util/strs"
)

type SearchOp struct {
	Field string
	Op    string
	Data  interface{}
}

func NewSearchParser(searchTag, scalarDefOp, fieldDefOp string, fieldDefSnakeCase ...bool) *SearchParser {
	searchTag = strings.TrimSpace(searchTag)
	scalarDefOp = strings.TrimSpace(scalarDefOp)
	fieldDefOp = strings.TrimSpace(fieldDefOp)
	if searchTag == "" || scalarDefOp == "" || fieldDefOp == "" {
		panic("params can't be empty")
	}
	sp := &SearchParser{
		searchTag:         "search",
		scalarDefOp:       scalarDefOp,
		sliceDefOp:        fieldDefOp,
		fieldDefSnakeCase: true,
		omitEmptyDef:      true,
		tokenOp:           "op",
		tokenAddOp:        "addOp",
		tokenOmitEmpty:    "omitEmpty",
		tokenField:        "field",
		AddOps:            make(map[string]func(interface{}) interface{}),
	}
	if len(fieldDefSnakeCase) > 0 {
		sp.fieldDefSnakeCase = fieldDefSnakeCase[0]
	}
	return sp
}

type SearchParser struct {
	searchTag         string
	scalarDefOp       string
	sliceDefOp        string
	fieldDefSnakeCase bool
	omitEmptyDef      bool
	tokenOp           string
	tokenAddOp        string
	tokenOmitEmpty    string
	tokenField        string
	AddOps            map[string]func(interface{}) interface{}
}

func (sp *SearchParser) SetSearchTag(tag string) *SearchParser {
	sp.searchTag = tag
	return sp
}

func (sp *SearchParser) SetTokenOp(token string) *SearchParser {
	sp.tokenOp = token
	return sp
}

func (sp *SearchParser) SetTokenAddOp(token string) *SearchParser {
	sp.tokenAddOp = token
	return sp
}

func (sp *SearchParser) SetTokenOmitEmpty(token string) *SearchParser {
	sp.tokenOmitEmpty = token
	return sp
}

func (sp *SearchParser) SetTokenField(token string) *SearchParser {
	sp.tokenField = token
	return sp
}

func (sp *SearchParser) SetOmitEmptyDef(omitEmpty bool) *SearchParser {
	sp.omitEmptyDef = omitEmpty
	return sp
}

func (sp *SearchParser) SetAddOpt(addOptName string, fn func(interface{}) interface{}) *SearchParser {
	sp.AddOps[addOptName] = fn
	return sp
}

func (sp *SearchParser) Parse(data interface{}) []*SearchOp {
	ret := make([]*SearchOp, 0)
	refData := ref.New(data)
	if !refData.IsStruct() {
		return ret
	}

	fieldsName, _ := refData.StructFieldsName()
	for _, fieldName := range fieldsName {
		fieldVal, _ := refData.StructValueGet(fieldName)
		if fieldVal == nil { // 字段值为空指针, 跳过
			continue
		}
		refFV := ref.New(fieldVal)
		if !refFV.IsStringOrNumber() || !refFV.IsSlice() { // 字段值不是字符串,数字,切片; 跳过
			continue
		}

		tag, ok := refData.StructTagGet(fieldName, sp.searchTag)
		if !ok { // 没有找到 查找的 标志tag
			continue
		}

		so := &SearchOp{
			Field: strs.SnakeString(fieldName),
			Data:  fieldVal,
		}
		if refFV.IsStringOrNumber() {
			so.Op = sp.scalarDefOp
		} else {
			so.Op = sp.sliceDefOp
		}

		addOps := make([]string, 0)
		omitEmpty := true

		for _, tmp := range strings.Split(tag, ",") {
			if strings.TrimSpace(tmp) == "" {
				continue
			}
			s := strings.SplitN(tmp, ":", 2)
			k, v := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
			if v == "" {
				so.Field = k
				continue
			}
			switch k {
			case sp.tokenField:
				so.Field = v
			case sp.tokenOp:
				so.Op = v
			case sp.tokenOmitEmpty:
				v := strings.ToUpper(v)
				if v == "F" || v == "FALSE" || v == "0" {
					omitEmpty = false
				}
			case sp.tokenAddOp:
				for _, o := range strings.Split(v, "|") {
					if ot := strings.TrimSpace(o); ot != "" {
						addOps = append(addOps, ot)
					}
				}
			}
		}

		if omitEmpty && refFV.GetValue().IsZero() {
			continue
		}

		if len(addOps) > 0 {
			for _, opName := range addOps {
				if fn, ok := sp.AddOps[opName]; ok {
					so.Data = fn(so.Data)
				}
			}
		}

		ret = append(ret, so)
	}

	return ret
}
