package strs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func IsEmpty(str string) bool {
	return str == ""
}

func StrWithFallback(str, fallback string) string {
	if !IsEmpty(str) {
		return str
	}

	return fallback
}

func ToInt64WithFallback(str string, fallback int64) int64 {
	output, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fallback
	}
	return output
}

func ToInt64WithDefaultZero(str string) int64 {
	return ToInt64WithFallback(str, 0)
}

func ToUint64WithFallback(str string, fallback uint64) uint64 {
	output, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fallback
	}
	return output
}

func ToUint64WithDefaultZero(str string) uint64 {
	return ToUint64WithFallback(str, 0)
}

func ToIntWithFallback(str string, fallback int) int {
	output, err := strconv.Atoi(str)
	if err != nil {
		return fallback
	}
	return output
}

func StrToIntWithDefaultZero(str string) int {
	return ToIntWithFallback(str, 0)
}

func ToFloat32WithFallback(str string, fallback float32) float32 {
	output, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return fallback
	}
	return float32(output)
}

func ToFloat32WithDefaultZero(str string) float32 {
	return ToFloat32WithFallback(str, 0)
}

func ToFloat64WithFallback(str string, fallback float64) float64 {
	output, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fallback
	}
	return output
}

func ToFloat64WithDefaultZero(str string) int {
	return ToIntWithFallback(str, 0)
}

func Stringify(v interface{}) string {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Array, reflect.Map:
		bs, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		return string(bs)
	default:
		return fmt.Sprint(v)
	}
}

/**
 * 驼峰转蛇形 snake string
 * @description XxYy to xx_yy , XxYY to xx_y_y
 * @date 2020/7/30
 * @param s 需要转换的字符串
 * @return string
 **/
func SnakeCase(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

/**
 * 蛇形转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY
 * @date 2020/7/30
 * @param s要转换的字符串
 * @return string
 **/
func CamelCase(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func SubString(s string, start, length int) string {
	if start < 0 || length <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) > start+length {
		return string(r[start : start+length])
	} else {
		return string(r[start:])
	}
}

func UcFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	data := []byte(s)
	if data[0] >= 'a' && data[0] <= 'z' {
		data[0] -= 32
	}

	return string(data)
}

func LcFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	data := []byte(s)
	if data[0] >= 'A' && data[0] <= 'Z' {
		data[0] += 32
	}

	return string(data)
}
