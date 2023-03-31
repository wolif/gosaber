package cooc

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

func isExported(name string) bool {
	// 是否包外部可访问(名字中首字母是否大写)
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func isExportedOrBuiltin(t reflect.Type) bool {
	// 循环剥离指针(解指针), 直到类型不是指针
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return isExported(t.Name()) || t.PkgPath() == "" // 包外部能访问 或者 类型所在的包名为空(golang内置类型)
}
