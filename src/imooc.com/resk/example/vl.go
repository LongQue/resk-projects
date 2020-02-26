package main

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Age       uint8  `validate:"gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

func main() {
	user :=&User{
		FirstName: "firstName",
		LastName:  "lastName",
		Age:       24,
		Email:     "zff@qq.com",
	}
	validate:=validator.New()
	//创建消息国际化通用翻译器
	cn:=zh.New()  //定义中文翻译器
	uni := ut.New(cn, cn)  //参数一缺省，参数二默认
	translator, found := uni.GetTranslator("zh")//返回默认，如果没有返回缺省，且是否找到默认
	if found{
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err!=nil{
			fmt.Println(err)
		}else {
			fmt.Println("not found")
		}
	}
	err := validate.Struct(user)
	if err!=nil{
		_,ok := err.(*validator.InvalidValidationError)
		if ok{
			fmt.Println(err)
		}
		errs,ok:=err.(validator.ValidationErrors)
		if ok{
			for _,err:=range errs{
				fmt.Println(err.Translate(translator))
			}
		}
	}

}