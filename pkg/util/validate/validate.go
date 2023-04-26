package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni               *ut.UniversalTranslator
	validate          *validator.Validate
	defTrans          ut.Translator
)

func init() {
	Init(
		map[string]locales.Translator{"zh": zh.New()},
		"zh",
		func(v *validator.Validate, trans ut.Translator) {
			zh_translations.RegisterDefaultTranslations(v, trans)
		},
	)
}

func Init(
	transes map[string]locales.Translator,
	defLang string,
	fnSetDefTrans func(validate *validator.Validate, defTrans ut.Translator),
) {
	validate = validator.New()
	transesTmp := []locales.Translator{}
	for _, tran := range transes {
		transesTmp = append(transesTmp, tran)
	}
	uni = ut.New(transesTmp[0], transesTmp[1:]...)
	trans, found := uni.GetTranslator(defLang)
	if !found {
		panic(fmt.Sprintf("validate translations language [%s] not found", defLang))
	}
	fnSetDefTrans(validate, trans)
	defTrans = trans
}

func Struct(i interface{}, lang ...string) string {
	trans := defTrans
	if len(lang) > 0 {
		if transTmp, found := uni.GetTranslator(lang[0]); found {
			trans = transTmp
		}
	}
	errs := []string{}
	if err := validate.Struct(i); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			errs = append(errs, fieldErr.Translate(trans))
		}
		return strings.Join(errs, ",")
	}
	return ""
}
