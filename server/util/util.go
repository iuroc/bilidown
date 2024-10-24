package util

import (
	"encoding/json"
	"net/http"
)

// 统一的 JSON 响应结构
type Res struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 发送响应
func (r Res) Write(w http.ResponseWriter) {
	bs, err := json.Marshal(r)
	if err != nil {
		w.Write([]byte(`{"success":false,"message":"系统错误","data":null}`))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
