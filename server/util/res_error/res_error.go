package res_error

import (
	"net/http"

	"github.com/iuroc/bilidown/server/util"
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
