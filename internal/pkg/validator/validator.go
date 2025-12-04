package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans ut.Translator
)

// Init 初始化验证器翻译器
func Init() error {
	// 获取 gin 的 validator 引擎
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取 json tag 的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		uni := ut.New(zhT, zhT)

		// this is usually know or extracted from http 'Accept-Language' header
		// also see uni.Import(...)
		var found bool
		trans, found = uni.GetTranslator("zh")
		if !found {
			return fmt.Errorf("translator not found")
		}

		// 注册翻译器
		if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("init validator failed")
}

// Translate 翻译错误信息
func Translate(err error) string {
	var errs []string
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			errs = append(errs, e.Translate(trans))
		}
		return strings.Join(errs, "; ")
	}
	return err.Error()
}
