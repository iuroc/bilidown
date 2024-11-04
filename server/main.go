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
	"strconv"
	"sync"
	"time"

	"github.com/getlantern/systray"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	systray.Run(onReady, nil)
}

func onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("Bilidown 视频解析器")
	systray.SetTooltip("Bilidown 视频解析器")

	if !CheckFfmpegInstalled() {
		log.Fatalln("请将 ffmpeg 安装到环境变量 PATH 中")
	}

	openBrowserItem := systray.AddMenuItem("打开应用", "打开应用")
	port := 8098
	u := fmt.Sprintf("http://127.0.0.1:%d", port)
	go func() {
		for {
			<-openBrowserItem.ClickedCh
			OpenBrowser(u)
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
		err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(port), nil)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	time.Sleep(time.Second)

	OpenBrowser(u)

	wg.Wait()
}

func OpenBrowser(u string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", u)
	case "darwin":
		cmd = exec.Command("open", u)
	case "linux":
		cmd = exec.Command("xdg-open", u)
	default:
		return fmt.Errorf("不支持的操作系统")
	}
	return cmd.Start()
}

func getIcon() []byte {
	// 读取 static/favicon.ico 文件
	data, err := os.ReadFile("static/favicon.ico")
	if err != nil {
		log.Fatal(err)
	}
	return data
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

func CheckFfmpegInstalled() bool {
	_, err := os.Stat("bin/ffmpeg.exe")
	if runtime.GOOS == "windows" && !os.IsNotExist(err) {
		return true
	}
	err = exec.Command("ffmpeg", "-version").Run()
	return err == nil
}
