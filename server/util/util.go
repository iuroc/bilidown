package util

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
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

func GetDB(path ...string) *sql.DB {
	pathStr := ""
	if len(path) == 0 {
		pathStr = "./data.db"
	} else {
		pathStr = path[0]
	}
	db, err := sql.Open("sqlite3", pathStr)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CheckBVID(bvid string) bool {
	if !regexp.MustCompile("^BV1[a-zA-Z0-9]+").MatchString(bvid) {
		return false
	}
	return true
}
