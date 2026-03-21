package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"bilidown/router"
	"bilidown/util"

	"github.com/getlantern/systray"
	_ "modernc.org/sqlite"
)

const (
	HTTP_PORT = 8098      // 限定 HTTP 服务器端口
	HTTP_HOST = ""        // 限定 HTTP 服务器主机
	VERSION   = "v2.0.15" // 软件版本号，将影响托盘标题显示
)

var urlLocal = fmt.Sprintf("http://127.0.0.1:%d", HTTP_PORT)
var urlLocalUnix = fmt.Sprintf("%s?___%d", urlLocal, time.Now().UnixMilli())

func main() {
	checkFFmpeg()
	// 启动托盘程序
	systray.Run(onReady, nil)
}

func onReady() {
	// 设置托盘图标
	setIcon()
	// 设置托盘标题
	setTitle()
	// 设置托盘菜单
	setMenuItem()
	// 初始化数据表
	mustInitTables()
	// 配置和启动 HTTP 服务器
	mustRunServer()
	// 调用默认浏览器访问端口
	time.Sleep(time.Millisecond * 1000)
	openBrowser(urlLocalUnix)
	// 保持运行
	select {}
}

// checkFFmpeg 检测 ffmpeg 的安装情况，如果未安装则打印提示信息。
func checkFFmpeg() {
	if _, err := util.GetFFmpegPath(); err != nil {
		fmt.Println("🚨 FFmpeg is missing. Install it from https://www.ffmpeg.org/download.html or place it in ./bin, then restart the application.")
		select {}
	}
}

// 配置和启动 HTTP 服务器
func mustRunServer() {
	// 前端打包文件
	http.Handle("/", http.FileServer(http.Dir("static")))
	// 后端接口服务
	http.Handle("/api/", http.StripPrefix("/api", router.API()))
	// 启动 HTTP 服务器
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", HTTP_HOST, HTTP_PORT), nil)
		if err != nil {
			log.Fatal("http.ListenAndServe:", err)
		}
	}()
}

// openBrowser 调用系统默认浏览器打开指定 URL
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		log.Printf("openBrowser: %v.", errors.New("unsupported operating system"))
	}
	if err := cmd.Start(); err != nil {
		log.Printf("openBrowser: %v.", err)
	}
	fmt.Printf("Opened in default browser: %s.\n", url)
}

// setIcon 设置托盘图标
func setIcon() {
	var path string
	if runtime.GOOS == "windows" {
		path = "static/favicon.ico"
	} else {
		path = "static/favicon-32x32.png"
	}
	systray.SetIcon(mustReadFile(path))
}

// mustReadFile 返回文件字节内容
func mustReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("os.ReadFile:", err)
	}
	return data
}

// setTitle 设置托盘标题和工具提示
func setTitle() {
	title := "Bilidown"
	tooltip := fmt.Sprintf("%s 视频解析器 %s (port:%d)", title, VERSION, HTTP_PORT)
	// only available on Mac and Windows.
	systray.SetTooltip(tooltip)
}

// setMenuItem 设置托盘菜单
func setMenuItem() {
	openBrowserItemText := fmt.Sprintf("打开主界面 (port:%d)", HTTP_PORT)
	openBrowserItem := systray.AddMenuItem(openBrowserItemText, openBrowserItemText)
	go func() {
		for {
			<-openBrowserItem.ClickedCh
			openBrowser(urlLocalUnix)
		}
	}()

	aboutItemText := "Github 项目主页"
	aboutItem := systray.AddMenuItem(aboutItemText, aboutItemText)
	go func() {
		for {
			<-aboutItem.ClickedCh
			openBrowser("https://github.com/iuroc/bilidown")
		}
	}()

	exitItemText := "退出应用"
	exitItem := systray.AddMenuItem(exitItemText, exitItemText)
	go func() {
		<-exitItem.ClickedCh
		log.Printf("Bilidown has exited.")
		systray.Quit()
	}()
}

// mustInitTables 初始化数据表
func mustInitTables() {
	db := util.MustGetDB()
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS "field" (
		"name" TEXT PRIMARY KEY NOT NULL,
		"value" TEXT
	)`); err != nil {
		log.Fatalln("create table field:", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS "log" (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"content" TEXT NOT NULL,
		"create_at" text NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		log.Fatalln("create table log:", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS "task" (
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
		"download_type" text NOT NULL DEFAULT 'merge',
		"create_at" text NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		log.Fatalln("create table task:", err)
	}

	if _, err := util.GetCurrentFolder(db); err != nil {
		log.Fatalln("util.GetCurrentFolder:", err)
	}

	if err := initHistoryTask(db); err != nil {
		log.Fatalln("initHistoryTask:", err)
	}

	// 添加可能缺失的列（用于数据库迁移）
	if err := addMissingColumns(db); err != nil {
		log.Fatalln("addMissingColumns:", err)
	}
}

// addMissingColumns 添加可能缺失的列（用于数据库迁移）
func addMissingColumns(db *sql.DB) error {
	// 检查download_type列是否存在，如果不存在则添加
	// SQLite没有直接的方法检查列是否存在，我们尝试添加列并忽略错误
	// 使用事务确保操作原子性
	util.SqliteLock.Lock()
	_, _ = db.Exec(`ALTER TABLE "task" ADD COLUMN "download_type" TEXT DEFAULT 'merge'`)
	// 将现有记录中的NULL值更新为默认值'merge'
	_, _ = db.Exec(`UPDATE "task" SET "download_type" = 'merge' WHERE "download_type" IS NULL`)
	util.SqliteLock.Unlock()

	// 忽略错误，因为列可能已经存在
	// SQLite错误码为1表示列已存在
	return nil
}

// initHistoryTask 将上一次程序运行时未完成的任务进度全部变为 error
func initHistoryTask(db *sql.DB) error {
	util.SqliteLock.Lock()
	_, err := db.Exec(`UPDATE "task" SET "status" = 'error' WHERE "status" IN ('waiting', 'running')`)
	util.SqliteLock.Unlock()
	return err
}
