package gormhelper

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wolif/gosaber/pkg/log"
	"github.com/wolif/gosaber/pkg/ref"
	"strings"
)

type ScopeOpt = string

const (
	OptEQ        ScopeOpt = "="
	OptNEQ       ScopeOpt = "!="
	OptGT        ScopeOpt = ">"
	OptGTE       ScopeOpt = ">="
	OptLT        ScopeOpt = "<"
	OptLTE       ScopeOpt = "<="
	OptIN        ScopeOpt = "IN"
	OptNIN       ScopeOpt = "NOT IN"
	OptLIKE      ScopeOpt = "LIKE"
	OptNLike     ScopeOpt = "NOT LIKE"
	OptISNULL    ScopeOpt = "IS NULL"
	OptISNOTNULL ScopeOpt = "IS NOT NULL"
)

type optFunc struct {
	Func                  func(opt string, field string, value interface{}) func(db *gorm.DB) *gorm.DB
	WithoutCalAddOptValue bool
}

var optMap = map[string]*optFunc{
	OptEQ:        {Func: optCompareRel},
	OptNEQ:       {Func: optCompareRel},
	OptGT:        {Func: optCompareRel},
	OptGTE:       {Func: optCompareRel},
	OptLT:        {Func: optCompareRel},
	OptLTE:       {Func: optCompareRel},
	OptIN:        {Func: optRangeRel},
	OptNIN:       {Func: optRangeRel},
	OptLIKE:      {Func: optLikeRel},
	OptNLike:     {Func: optLikeRel},
	OptISNULL:    {Func: optNullRel, WithoutCalAddOptValue: true},
	OptISNOTNULL: {Func: optNullRel, WithoutCalAddOptValue: true},
}

func OptExist(opt ScopeOpt) bool {
	_, found := optMap[opt]
	return found
}

func ExtendOpt(opt string, withoutCalAddOptValue bool, fn func(opt ScopeOpt, field string, value interface{}) func(db *gorm.DB) *gorm.DB, options ...bool) {
	if opt == "" {
		return
	}
	recoverFn := true
	if len(options) > 0 {
		recoverFn = options[0]
	}

	if !OptExist(opt) || recoverFn {
		optMap[opt] = &optFunc{
			Func:                  fn,
			WithoutCalAddOptValue: withoutCalAddOptValue,
		}
	}
}

func resolveOmitEmptyTag(tag string) bool {
	return ref.IsInScope(strings.ToUpper(strings.TrimSpace(tag)), "", "Y", "YES", "1")
}

func where(field string, value interface{}, needOmitEmpty string, additionOpt addOpt) func(db *gorm.DB) *gorm.DB {
	if ref.New(value).IsSlice() {
		return whereOpt(field, OptIN, value, needOmitEmpty, additionOpt)
	}
	return whereOpt(field, OptEQ, value, needOmitEmpty, additionOpt)
}

func whereOpt(field string, opt string, value interface{}, needOmitEmpty string, addOpt addOpt) func(db *gorm.DB) *gorm.DB {
	opt = strings.TrimSpace(opt)
	field = strings.TrimSpace(field)
	if field == "" || opt == "" {
		return scopeNullFunc
	}

	if !OptExist(opt) {
		log.Errorf("opt [%s] gorm helper scope not support", opt)
		return scopeNullFunc
	}

	optFn, _ := optMap[opt]
	if optFn.WithoutCalAddOptValue {
		return optFn.Func(opt, field, value)
	}

	vRef := ref.New(value)
	if resolveOmitEmptyTag(needOmitEmpty) && vRef.IsZero() {
		return scopeNullFunc
	}

	// 数据额外操作
	if addOpt != "" {
		value = calAddOptValue(value, addOpt)
	}

	return optFn.Func(opt, field, value)
}

func optNullRel(opt string, field string, _ interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` %s", field, opt))
	}
}

func optLikeRel(opt string, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
	vRef := ref.New(value)
	if vRef.IsString() {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s` %s ?", field, opt), fmt.Sprintf("%%%s%%", value))
		}
	}
	return scopeNullFunc
}

func optRangeRel(opt string, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
	vRef := ref.New(value)
	if vRef.IsSlice() && vRef.GetValue().Len() > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s` %s (?)", field, opt), value)
		}
	}
	return scopeNullFunc
}

func optCompareRel(opt string, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
	vRef := ref.New(value)
	if vRef.IsString() || vRef.IsNumber() {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s` %s ?", field, opt), value)
		}
	}
	return scopeNullFunc
}
