package main

import (
	"bilidown/router"
	"bilidown/task"
	"bilidown/util"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if !CheckFfmpegInstalled() {
		log.Fatalln("请将 ffmpeg 安装到环境变量 PATH 中")
	}

	InitTables()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/api/", http.StripPrefix("/api", router.API()))

	fmt.Println("http://127.0.0.1:8098")
	http.ListenAndServe(":8098", nil)
}

// InitTables 初始化数据表
func InitTables() *sql.DB {
	db := util.GetDB()
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "field" (
		"name" TEXT PRIMARY KEY NOT NULL,
		"value" TEXT
	)`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "log" (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"content" TEXT NOT NULL,
		"create_at" text NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "task" (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"bvid" text NOT NULL,
		"cid" integer NOT NULL,
		"format" integer NOT NULL,
		"title" text NOT NULL,
		"owner" text NOT NULL,
		"cover" text NOT NULL,
		"status" text NOT NULL,
		"folder" text NOT NULL,
		"create_at" text NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = task.GetCurrentFolder(db)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CheckFfmpegInstalled() bool {
	err := exec.Command("ffmpeg", "-version").Run()
	return err == nil
}
