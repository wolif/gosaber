package structtags

import (
	"fmt"
	"github.com/wolif/gosaber/pkg/ref"
	"reflect"
	"strings"
	"sync"
)

// map[struct type][tag name][struct fieldName][tag fieldName][tag fieldValue]
var cache = make(map[reflect.Type]map[string]map[string]map[string]string)

var mu sync.Mutex

// 解析结构的中所有字段的某种tag
// exp. `tagName:"valWithoutKey,a:1,b:2,c:3"` => map[string][string]string{"id":map[string]string{"":"valWithoutKey","a":"1","b":"2","c":"3"}}
// 如果tag不存在,不会填充 defResult
// 特殊值,key为空字符串的值,只能有一个, 多个的话,后面的覆盖前面的
func Parse(structObj interface{}, tagName string, defResults ...map[string]string) (map[string]map[string]string, error) {
	refObj := ref.New(structObj)
	if !refObj.IsStruct() {
		return nil, fmt.Errorf("only struct kind data can be parsed here")
	}
	tagName = strings.TrimSpace(tagName)
	refObjType := refObj.GetType()
	if tags, ok := cache[refObjType]; ok {
		if ret, ok := tags[tagName]; ok {
			return ret, nil
		}
	}

	mu.Lock()
	defer mu.Unlock()

	defResult := make(map[string]string)
	if len(defResults) > 0 {
		defResult = defResults[0]
	}

	if _, ok := cache[refObjType]; !ok {
		cache[refObjType] = make(map[string]map[string]map[string]string)
	}

	fieldNames, _ := refObj.StructFieldsName()
	for _, fieldName := range fieldNames {
		oneFieldRes := make(map[string]string)
		tagStr, ok := refObj.StructTagGet(fieldName, tagName)
		if !ok { // tag 不存在
			continue
		}
		if _, ok := cache[refObjType][tagName]; !ok {
			cache[refObjType][tagName] = make(map[string]map[string]string)
		}
		for _, kv := range strings.Split(strings.TrimSpace(tagStr), ",") {
			if strings.TrimSpace(kv) == "" {
				continue
			}
			kvSlice := strings.SplitN(kv, ":", 2)
			if len(kvSlice) == 1 { // 无key值
				if len(strings.TrimSpace(kvSlice[0])) > 0 {
					oneFieldRes[""] = strings.TrimSpace(kvSlice[0])
				}
			} else if tagFiledName := strings.TrimSpace(kvSlice[0]); tagFiledName != "" {
				oneFieldRes[tagFiledName] = strings.TrimSpace(kvSlice[1])
			}
		}
		for dk, dv := range defResult {
			if tagFieldVal, ok := oneFieldRes[dk]; !ok || tagFieldVal == "" {
				oneFieldRes[dk] = dv
			}
		}
		cache[refObjType][tagName][fieldName] = oneFieldRes
	}
	return cache[refObjType][tagName], nil
}
