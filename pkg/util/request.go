package util

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// 请求参数
func getRequestBody(context *gin.Context, s interface{}) error { //获取request的body
	body, _ := context.Get("json") //转换成json格式
	reqBody, _ := body.(string)
	decoder := json.NewDecoder(bytes.NewReader([]byte(reqBody)))
	decoder.UseNumber() //作为数字而不是float64
	err := decoder.Decode(&s)//从body中获取的参数存入s中
	return err
}

// GetPostParams 获取post接口参数
func GetPostParams(ctx *gin.Context) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	err := getRequestBody(ctx, &params)
	return params, err
}

// GetUserInfo 获取用户信息
func GetUserInfo(ctx *gin.Context) interface{} {
	example := ctx.MustGet("userInfo")
	return example
}