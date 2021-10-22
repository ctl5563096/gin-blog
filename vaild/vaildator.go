package valid

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)


var (
	uni      *ut.UniversalTranslator
	trans ut.Translator
)

func Init() {
	//注册翻译器
	zhTr := zh.New()
	uni = ut.New(zhTr, zhTr)

	trans, _ = uni.GetTranslator("zh")

	//获取gin的校验器
	validate := binding.Validator.Engine().(*validator.Validate)
	//注册翻译器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return 
	}
}

//Translate 翻译错误信息
func Translate(err error) map[string][]string {

	var result = make(map[string][]string)

	errors := err.(validator.ValidationErrors)

	for _, err := range errors{
		result[err.Field()] = append(result[err.Field()], err.Translate(trans))
	}
	return result
}