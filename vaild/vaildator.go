package valid

import (
	"gin-blog/pkg/util"
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

type ValidateError struct {
	s string
}

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

func (fe *ValidateError) Error() string {
	return fe.s
}

// RequestData 验证错误
func RequestData(params interface{}) error {
	var errorMap map[string][]string
	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = Translate(err)
			//循环遍历Map 只返回第一个错误信息
			for _,v:= range errorMap{
				for _,z := range v{
					util.WriteLog("business_error",4,z)
					return &ValidateError{z}
				}
			}
		default:
			return &ValidateError{"未知错误!"}
		}
	}
	return nil
}