package router

import (
	"bilidown/task"
	"bilidown/util"
	"encoding/json"
	"net/http"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		util.Res{Success: false, Message: "不支持的请求方法"}.Write(w)
		return
	}
	var body task.TaskInitOption
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	if !util.CheckBVID(body.Bvid) {
		util.Res{Success: false, Message: "bvid 格式错误"}.Write(w)
		return
	}
	if body.Cover == "" || body.Title == "" || body.Owner == "" {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
	}

	if !util.IsValidURL(body.Cover) {
		util.Res{Success: false, Message: "封面链接格式错误"}.Write(w)
		return
	}
	if !util.IsValidFormatCode(body.Format) {
		util.Res{Success: false, Message: "清晰度代码错误"}.Write(w)
		return
	}
	db := util.GetDB()
	body.Folder, err = task.GetCurrentFolder(db)
	body.Status = "waiting"
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	_task := task.Task{TaskInitOption: body}
	err = _task.Create(db)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "创建成功", Data: struct {
		ID int64 `json:"id"`
	}{ID: _task.ID}}.Write(w)
	go _task.Start(db)
}
