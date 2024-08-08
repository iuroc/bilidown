package main

import (
	"fmt"

	"github.com/iuroc/bilidown"
)

func main() {
	ClearTerminal()
	cookieValue := promptLogin()
	ClearTerminal()
	promptDownload(cookieValue)
}

func promptDownload(cookieValue string) {
	for {
		fmt.Print("> 请输入 Bilibili 视频链接：")
		var url string
		fmt.Scan(&url)
		videoId, err := bilidown.CheckVideoURLOrID(url)
		if err != nil {
			ClearTerminal()
			fmt.Print("❗️ 您输入的视频链接格式错误，请重新输入\n\n")
			promptDownload(cookieValue)
			continue
		}
		videoURL := bilidown.MakeVideoURL(videoId)
		bilidown.ParseVideo(videoURL, cookieValue)
	}
}

// promptLogin 首先检查本地 Cookie，如果无可用 Cookie，则通过 Select 让用户选择是否登录，
// 如果用户选择登录，则调用浏览器进行登录，并保存返回的 Cookie，否则 Cookie 保持空值表示游客访问。
func promptLogin() (cookieValue string) {
	cookieSavePath := "cookie"
	cookieValue, err := bilidown.GetCookieValue(cookieSavePath)
	if err != nil {
		for {
			if shouldLogin() {
				cookie, err := bilidown.Login()
				if err != nil {
					continue
				}
				bilidown.SaveCookie(cookie, cookieSavePath)
				cookieValue = cookie.Value
			}
			break
		}
	}
	return cookieValue
}

func shouldLogin() bool {
	items := []string{"登录账号（支持全部分辨率）", "游客访问（仅支持低分辨率）"}
	fmt.Println("🔅 当前未登录，请选择是否登录：")
	for index, item := range items {
		fmt.Printf("  %d. %s\n", index+1, item)
	}
	fmt.Printf("> 请输入操作序号 [%d-%d]：", 1, len(items))
	var id int
	_, err := fmt.Scanf("%d\n", &id)
	if err != nil || id > len(items) {
		ClearTerminal()
		fmt.Print("❗️ 您输入的序号错误，请重新输入\n\n")
		return shouldLogin()
	}
	return id == 1
}

func ClearTerminal() {
	fmt.Print("\x1b[H\x1b[2J")
}
