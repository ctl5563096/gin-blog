package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// GetTracingSpan 获取调用链路父级初始化
func GetTracingSpan(c *gin.Context) opentracing.SpanContext {
	parentSpanContext, isExists := c.Get("ParentSpanContext")
	if isExists {
		return parentSpanContext.(opentracing.SpanContext)
	} else {
		// 容错处理！
		return opentracing.StartSpan("ParentSpanContextError").Context()
	}
}

// ParseQueryRequest Query参数验证
func ParseQueryRequest(c *gin.Context, request interface{}) (err error) {
	err = c.ShouldBindQuery(request)
	handleParamsError(c, err)
	return
}

// ParseRequest 普通参数验证
func ParseRequest(c *gin.Context, request interface{}) (err error) {
	err = c.ShouldBind(request)
	handleParamsError(c, err)
	return
}

// 处理422错误信息
func handleParamsError(c *gin.Context, err error) {
	if err != nil {
		var errStr string
		switch err.(type) {
			case validator.ValidationErrors:
				errStr = err.Error()
			case *json.UnmarshalTypeError:
				unmarshalTypeError := err.(*json.UnmarshalTypeError)
				errStr = fmt.Errorf("%s [类型错误，期望类型] %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
			default:
				// logging.Error("unknown error：" + err.Error())
				errStr = errors.New("参数校验[unknown error]").Error()
		}
		// 此处错误信息该修正下
		FailWithParameter(errStr, c)
	}
}