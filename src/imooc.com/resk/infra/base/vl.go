package base

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	vzth "gopkg.in/go-playground/validator.v9/translations/zh"
	"resk-projects/infra"
)

var validate *validator.Validate
var translator ut.Translator

func Validate() *validator.Validate {
	return validate
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func Transtate() ut.Translator {
	Check(translator)
	return translator
}
func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
	//创建消息国际化通用翻译器
	cn := zh.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		err := vzth.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Error("Not found translate")
	}
}

func ValidateStruct(s interface{}) (err error) {
	//验证
	err = Validate().Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logrus.Error(err)
		}

		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				logrus.Error(err.Translate(translator))
			}
		}
		return err
	}
	return nil
}
