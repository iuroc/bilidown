package res_error

import (
	"net/http"
	"strings"

	"bilidown/util"
)

func sendError(w http.ResponseWriter, message string) {
	util.Res{Message: message, Success: false}.Write(w)
}

func BvidFormatError(w http.ResponseWriter) {
	sendError(w, "BVID 格式错误")
}

func URLFormatError(w http.ResponseWriter) {
	sendError(w, "URL 格式错误")
}

func ParamError(w http.ResponseWriter) {
	sendError(w, "参数错误")
}

func MethodNotAllow(w http.ResponseWriter) {
	sendError(w, "不允许的请求类型")
}

func NoRedirectedLocation(w http.ResponseWriter) {
	sendError(w, "未发现重定向目标")
}

func FileNotExist(w http.ResponseWriter, path ...string) {
	message := "文件不存在"
	if len(path) > 0 {
		message += ": " + strings.Join(path, ", ")
	}
	sendError(w, message)
}

func SystemError(w http.ResponseWriter) {
	sendError(w, "系统错误")
}

func FileTypeNotAllow(w http.ResponseWriter) {
	sendError(w, "不支持的文件类型")
}
