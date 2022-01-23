package app

import (
	"gin-blog/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Result Response setting gin.JSON
func Result(httpCode, errCode int, errMsg string, data interface{}, c *gin.Context) {
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  errMsg,
		Data: data,
	})
}

// OK Response Success
func OK(c *gin.Context) {
	Result(http.StatusOK, e.SUCCESS, "操作成功", map[string]interface{}{}, c)
}

// OkWithMessage Response OkWithData
func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, e.SUCCESS, message, map[string]interface{}{}, c)
}

// OkWithData Response OkWithData
func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, e.SUCCESS, "OK", data, c)
}

// OkWithCodeData Response OkWithCodeData
func OkWithCodeData(message string, data interface{}, code int, c *gin.Context) {
	Result(http.StatusOK, code, message, data, c)
}

// Fail Response Fail
func Fail(c *gin.Context) {
	Result(http.StatusOK, e.ERROR, "操作失败", map[string]interface{}{}, c)
}

// FailWithMessage Response FailWithMessage
func FailWithMessage(message string, code int, c *gin.Context) {
	Result(http.StatusOK, code, message, map[string]interface{}{}, c)
}

// FailWithParameter Response FailWithMessage
func FailWithParameter(message string, c *gin.Context) {
	Result(http.StatusOK, e.INVALID_PARAMS, message, map[string]interface{}{}, c)
}

// LoseWithParameter  Response LoseWithParameter
func LoseWithParameter(message string, c *gin.Context) {
	Result(http.StatusOK, e.MISS_PARAMS, message, map[string]interface{}{}, c)
}

// MissToken token验证失败
func MissToken(c *gin.Context)  {
	Result(e.MISS_TOKEN,e.MISS_TOKEN,e.GetMsg(e.MISS_TOKEN),map[string]interface{}{},c)
}