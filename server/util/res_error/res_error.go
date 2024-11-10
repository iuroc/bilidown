package res_error

import (
	"bilidown/server/util"
	"net/http"
)

func SendError(w http.ResponseWriter, message string) {
	util.Res{Message: message, Success: false}.Write(w)
}

func BvidFormatError(w http.ResponseWriter) {
	SendError(w, "BVID 格式错误")
}

func ParamError(w http.ResponseWriter) {
	SendError(w, "参数错误")
}

func MethodNotAllow(w http.ResponseWriter) {
	SendError(w, "不允许的请求类型")
}
