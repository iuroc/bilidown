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
	var body []task.TaskInitOption
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	db := util.GetDB()
	defer db.Close()
	for _, item := range body {
		if !util.CheckBVID(item.Bvid) {
			util.Res{Success: false, Message: "bvid 格式错误"}.Write(w)
			return
		}
		if item.Cover == "" || item.Title == "" || item.Owner == "" {
			util.Res{Success: false, Message: "参数错误"}.Write(w)
		}

		if !util.IsValidURL(item.Cover) {
			util.Res{Success: false, Message: "封面链接格式错误"}.Write(w)
			return
		}
		if !util.IsValidFormatCode(item.Format) {
			util.Res{Success: false, Message: "清晰度代码错误"}.Write(w)
			return
		}
		item.Folder, err = task.GetCurrentFolder(db)
		item.Status = "waiting"
		if err != nil {
			util.Res{Success: false, Message: err.Error()}.Write(w)
			return
		}
		_task := task.Task{TaskInitOption: item}
		err = _task.Create(db)
		if err != nil {
			util.Res{Success: false, Message: err.Error()}.Write(w)
			return
		}
		go _task.Start()
	}
	util.Res{Success: true, Message: "创建成功"}.Write(w)
}

func GetActiveTask(w http.ResponseWriter, r *http.Request) {
	util.Res{Success: true, Data: struct {
		Tasks []*task.Task `json:"tasks"`
	}{
		Tasks: task.GlobalTaskList,
	}}.Write(w)
}
