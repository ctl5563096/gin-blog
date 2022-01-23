package e

var errorMsg = map[int]string{
	SUCCESS 							: "ok",
	ERROR 								: "fail",
	INVALID_PARAMS 						: "请求参数错误",
	ERROR_EXIST_TAG 					: "已存在该标签名称",
	ERROR_NOT_EXIST_TAG 				: "该标签不存在",
	ERROR_NOT_EXIST_ARTICLE 			: "该文章不存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL 		: "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT 		: "Token已超时",
	ERROR_AUTH_TOKEN 					: "Token生成失败",
	ERROR_AUTH 							: "Token错误",
	MISS_PARAMS 						: "参数缺失,请验证参数",
	PASSWORD_ERROR 						: "密码错误,请重新登陆",
	MISS_TOKEN							: "缺少token,请重新检查",
	USER_IN_BLACK						: "用户在黑名单里面,无法使用",
	SYSTEM_ERROR						: "系统错误",
	USER_NOT_EXIST						: "用户不存在",
	UPDATE_ERROR						: "更新失败",
	INSERT_ERROR						: "插入失败",
	DELETE_ERROR						: "删除失败",
	OPERATION_ONT_PERMITTED				: "操作不被允许",
	IMAGE_TYPE_ERROR					: "上传图片类型不被允许",
	UPLOAD_FAIL							: "上传失败",
	FIND_ARTICLE_ERROR					: "搜索文章失败",
	PARAMS_ERROR						: "参数错误",
	FILE_TYPE_ERROR						: "文件格式错误,请检查文件格式是否符合:",
	CREATE_THUMB_ERROR				 	: "生成缩略图失败",
}

func GetMsg(code int) string {
	msg,ok := errorMsg[code]
	if ok {
		return msg
	}
	return errorMsg[ERROR]
}
