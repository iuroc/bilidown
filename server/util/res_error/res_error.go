package res_error

import (
	"net/http"

	"bilidown/util"
)

// Send 发送异常响应
func Send(w http.ResponseWriter, message string) {
	util.Res{Message: message, Success: false}.Write(w)
}

const (
	BvidFormatError       = "错误的 Bvid 格式"
	URLFormatError        = "错误的 URL 格式"
	MidFormatError        = "错误的 Mid 格式"
	SeasonIdFormatError   = "错误的 SeasonId 格式"
	ParamError            = "参数错误"
	MethodNotAllowError   = "不允许的请求方式"
	NoLocationError       = "无重定向目标地址"
	FileNotFountError     = "文件不存在"
	FileTypeNotAllowError = "不允许的文件类型"
	SystemError           = "系统错误"
	NotLogin              = "未登录"
)
