package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/iuroc/bilidown"
	"github.com/iuroc/gododo/biliqr"
	"github.com/skip2/go-qrcode"
)

var scanner = bufio.NewScanner(os.Stdin)
var loginQrcodeUrl *qrcode.QRCode

const qrcodeUrl = "http://127.0.0.1:16666/login"
const serverAddr = ":16666"

func main() {
	bilidown.InitDir("download")
	bilidown.InitDir("temp")
	go InitHttpService()
	cookieValue := promptLogin()
	ClearTerminal()
	promptDownload(cookieValue)
}

func RequireQR() (string, error) {
	ClearTerminal()
	fmt.Println("正在获取登录二维码...")
	qr, info, err := biliqr.NewLoginQR(qrcode.Low)
	if err != nil {
		log.Fatalln(err.Error())
	}
	loginQrcodeUrl = qr
	ClearTerminal()
	bilidown.OpenUrlByBrowser(qrcodeUrl)
	for {
		status, err := biliqr.GetQRStatus(info.OauthKey)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if status.Code == 0 {
			return status.SESSDATA, nil
		}
		if status.Code == -2 {
			return "", errors.New("二维码已过期")
		}
	}
}

func promptDownload(cookieValue string) {
	for {
		fmt.Print("> 请输入 Bilibili 视频链接: ")
		if !scanner.Scan() {
			log.Fatal(scanner.Err())
		}
		url := scanner.Text()
		videoId, err := bilidown.CheckVideoURLOrID(url)
		if err != nil {
			ClearTerminal()
			fmt.Print("❗️ 您输入的视频链接格式错误, 请重新输入\n\n")
			continue
		}
		videoURL := bilidown.MakeVideoURL(videoId)
		parseResult, err := bilidown.ParseVideo(videoURL, cookieValue)
		if err != nil {
			ClearTerminal()
			fmt.Println(err.Error())
			fmt.Print("❗️ 视频解析失败, 请重试\n\n")
			continue
		}
		// https://www.bilibili.com/video/BV1fK4y1t7hj/
		fmt.Println()
		fmt.Print("👇👇👇👇👇👇👇👇 解析成功, 以下是解析结果 👇👇👇👇👇👇👇👇\n\n")
		fmt.Printf("🌟 视频标题: %s\n📝 视频描述: %s\n\n", parseResult.Title, parseResult.Desc)
		fmt.Print("🔸🔸🔸🔸🔸🔸🔸🔸🔸🔸 视频信息 🔸🔸🔸🔸🔸🔸🔸🔸🔸🔸\n")
		for index, staff := range parseResult.Staff {
			fmt.Printf("🔹 %s: %s", staff.Title, staff.Name)
			if index%3 == 2 || index == len(parseResult.Staff)-1 {
				fmt.Println()
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Printf("👍 点赞: %d  💰 投币: %d  ⭐ 收藏: %d  🔗 分享: %d\n🎬 播放: %d  💬 弹幕: %d  📝 评论: %d  🕒 发布: %s\n\n",
			parseResult.Stat.Like,
			parseResult.Stat.Coin,
			parseResult.Stat.Favorite,
			parseResult.Stat.Share,
			parseResult.Stat.View,
			parseResult.Stat.Danmaku,
			parseResult.Stat.Reply,
			parseResult.PubdateString(),
		)
		fmt.Print("🔸🔸🔸🔸🔸🔸🔸🔸🔸🔸 下载选项 🔸🔸🔸🔸🔸🔸🔸🔸🔸🔸\n")
		for index, video := range parseResult.Dash.Video {
			fmt.Printf("[%d]\t[分辨率: %s]\t[比特率: %dKbps]\t[编解码: %s]\n", index+1, video.Description(parseResult.SupportFormats), int(video.Bandwidth/1000), video.Codecs)
		}
		fmt.Println()
		var videoSelectNum int
		// https://www.bilibili.com/video/BV1fK4y1t7hj/
		cancel := false
		for {
			fmt.Printf("请输入需要下载的视频序号 [%d-%d], 输入 0 取消当前操作: ", 1, len(parseResult.Dash.Video))
			if !scanner.Scan() {
				log.Fatal(scanner.Err())
			}
			videoSelectNum, err = strconv.Atoi(scanner.Text())
			if err != nil || videoSelectNum < 0 || videoSelectNum > len(parseResult.Dash.Video) {
				fmt.Print("❗️ 请输入正确的序号\n\n")
				continue
			}
			if videoSelectNum == 0 {
				cancel = true
			}
			break
		}
		ClearTerminal()
		if cancel {
			continue
		}
		outputPath, err := bilidown.Download(parseResult, videoSelectNum-1, "download", "temp")
		if err != nil {
			ClearTerminal()
			fmt.Printf("❗️ 视频下载失败, 建议您稍后重试: %v\n\n", err)
			continue
		}
		absPath, err := filepath.Abs(outputPath)
		if err != nil {
			log.Fatalln(absPath)
		}
		ClearTerminal()
		fmt.Printf("视频下载成功: %s\n\n", absPath)
		fmt.Println("🚗 回车继续解析下一个视频")
		if !scanner.Scan() {
			log.Fatalln(scanner.Err())
		}
		ClearTerminal()
	}
}

// promptLogin 首先检查本地 Cookie, 如果无可用 Cookie, 则通过 Select 让用户选择是否登录,
// 如果用户选择登录, 则调用浏览器进行登录, 并保存返回的 Cookie, 否则 Cookie 保持空值表示游客访问。
func promptLogin() (cookieValue string) {
	cookieSavePath := "cookie"
	cookieValue, err := bilidown.GetCookieValue(cookieSavePath)
	if err != nil {
		for {
			if shouldLogin() {
				cookieValue, err = RequireQR()
				if err != nil {
					ClearTerminal()
					if err.Error() != "context canceled" {
						fmt.Print("❗️ " + err.Error() + "\n\n")
					}
					continue
				}
				bilidown.SaveCookie(&network.Cookie{
					Value:   cookieValue,
					Name:    "SESSDATA",
					Expires: float64(time.Now().Add(160 * 24 * time.Hour).Unix()),
				}, cookieSavePath)
			}
			break
		}
	}
	return cookieValue
}

// shouldLogin 返回是否应该调用浏览器进行登录操作
func shouldLogin() bool {
	items := []string{"登录账号（支持全部分辨率）", "游客访问（仅支持低分辨率）"}
	fmt.Println("🔅 当前未登录, 请选择是否登录: ")
	for index, item := range items {
		fmt.Printf("  %d. %s\n", index+1, item)
	}
	fmt.Printf("> 请输入操作序号 [%d-%d]: ", 1, len(items))
	if !scanner.Scan() {
		log.Fatal(scanner.Err())
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil || id <= 0 || id > len(items) {
		ClearTerminal()
		fmt.Print("❗️ 您输入的序号错误, 请重新输入\n\n")
		return shouldLogin()
	}
	return id == 1
}

func ClearTerminal() {
	fmt.Print("\x1b[H\x1b[2J")
}

func HandleLoginQrcode(w http.ResponseWriter, r *http.Request) {
	if loginQrcodeUrl == nil {
		fmt.Fprintf(w, "二维码未获取，请获取后访问")
		return
	}
	loginQrcodeUrl.WriteFile(100, "1111.png")
	qrcodeImg, err := loginQrcodeUrl.PNG(300)
	if err != nil {
		fmt.Fprintf(w, "二维码生成失败"+err.Error())
		return
	}
	w.Write(qrcodeImg)
	return
}

func InitHttpService() {
	http.HandleFunc("/login", HandleLoginQrcode)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
