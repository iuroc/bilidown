package router

import (
	"bilidown/util"
	"encoding/json"
	"net/http"
)

type RequestData struct {
	Bvid   string `json:"bvid"`
	Cid    int    `json:"cid"`
	Format int    `json:"format"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		util.Res{Success: false, Message: "不支持的请求方法"}.Write(w)
	}
	var body RequestData
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	if !util.CheckBVID(body.Bvid) {
		util.Res{Success: false, Message: "bvid 格式错误"}.Write(w)
		return
	}
	// db := util.GetDB()

	// TODO
}
