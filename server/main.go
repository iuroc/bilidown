package main

import (
	"bilidown/router"
	"bilidown/task"
	"bilidown/util"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/getlantern/systray"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if _, err := util.GetFFmpegPath(); err != nil {
		fmt.Println("🚨 FFmpeg is missing. Install it from https://www.ffmpeg.org/download.html or place it in ./bin, then restart the application.")
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()
	}
	systray.Run(onReady, nil)
}

const HTTP_PORT = 8098
const HTTP_HOST = "127.0.0.1"

func onReady() {
	if icon, err := getIcon(); err != nil {
		log.Fatalln(err)
	} else {
		systray.SetIcon(icon)
	}

	systray.SetTitle("Bilidown 视频解析器")
	systray.SetTooltip("Bilidown 视频解析器")

	_url := fmt.Sprintf("http://%s:%d", HTTP_HOST, HTTP_PORT)

	openBrowserItem := systray.AddMenuItem("打开应用", "打开应用")
	go func() {
		for {
			<-openBrowserItem.ClickedCh
			OpenBrowser(fmt.Sprintf("%s?_=%d", _url, time.Now().UnixNano()))
		}
	}()

	exitItem := systray.AddMenuItem("退出应用", "退出应用")
	go func() {
		<-exitItem.ClickedCh
		systray.Quit()
	}()

	db := util.GetDB()
	InitTables(db)
	task.InitHistoryTask(db)
	db.Close()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/api/", http.StripPrefix("/api", router.API()))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", HTTP_HOST, HTTP_PORT), nil)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	time.Sleep(time.Millisecond * 1000)

	OpenBrowser(fmt.Sprintf("%s?_=%d", _url, time.Now().UnixNano()))

	wg.Wait()
}

func OpenBrowser(_url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", _url)
	case "darwin":
		cmd = exec.Command("open", _url)
	case "linux":
		cmd = exec.Command("xdg-open", _url)
	default:
		return fmt.Errorf("不支持的操作系统")
	}
	return cmd.Start()
}

func getIcon() ([]byte, error) {
	// 读取 static/favicon.ico 文件
	return os.ReadFile("static/favicon.ico")
}

// InitTables 初始化数据表
func InitTables(db *sql.DB) {
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
		"duration" integer NOT NULL,
		"create_at" text NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = task.GetCurrentFolder(db)
	if err != nil {
		log.Fatalln(err)
	}
}
