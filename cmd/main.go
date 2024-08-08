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
		fmt.Print("> 请输入 Bilibili 视频链接: ")
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
		parseResult, err := bilidown.ParseVideo(videoURL, cookieValue)
		if err != nil {
			ClearTerminal()
			fmt.Print("❗️ 视频解析失败，请重试\n\n")
			promptDownload(cookieValue)
			continue
		}
		// https://www.bilibili.com/video/BV1fK4y1t7hj/
		ClearTerminal()
		fmt.Print("👇👇👇👇👇👇👇👇 解析成功，以下是解析结果 👇👇👇👇👇👇👇👇\n\n")
		fmt.Printf("🌟 视频标题: %s\n📝 视频描述: %s\n\n", parseResult.Title, parseResult.Desc)
		fmt.Print("🔸🔸🔸🔸🔸🔸 视频信息 🔸🔸🔸🔸🔸🔸\n")
		for index, staff := range parseResult.Staff {
			fmt.Printf("🔹 %s: %s ", staff.Title, staff.Name)
			if index%3 == 2 || index == len(parseResult.Staff)-1 {
				fmt.Println()
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Printf("👍 点赞: %d 💰 投币: %d ⭐ 收藏: %d 🔗 分享: %d\n🎬 播放: %d 💬 弹幕: %d 📝 评论: %d 🕒 发布: %s\n\n",
			parseResult.Stat.Like,
			parseResult.Stat.Coin,
			parseResult.Stat.Favorite,
			parseResult.Stat.Share,
			parseResult.Stat.View,
			parseResult.Stat.Danmaku,
			parseResult.Stat.Reply,
			parseResult.PubdateString(),
		)
		fmt.Print("🔸🔸🔸🔸🔸🔸 下载选项 🔸🔸🔸🔸🔸🔸\n")
		for index, video := range parseResult.Dash.Video {
			fmt.Printf("[%d]\t[%s]\t[%dKbps]\t[%s]\n", index+1, video.Description(parseResult.SupportFormats), int(video.Bandwidth/1000), video.Codecs)
		}
		fmt.Printf("\n\n请输入需要下载的视频序号 [%d-%d]：", 1, len(parseResult.Dash.Video))
		var videoSelectIndex int
		fmt.Scan(&videoSelectIndex)
		fmt.Println("🚗 回车继续解析下一个视频")
		fmt.Scanln()
		fmt.Scanln()
		ClearTerminal()
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
					ClearTerminal()
					if err.Error() != "context canceled" {
						fmt.Print("❗️ 打开浏览器失败，请确保安装了 Chrome 浏览器\n\n")
					}
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

// shouldLogin 返回是否应该调用浏览器进行登录操作
func shouldLogin() bool {
	items := []string{"登录账号（支持全部分辨率）", "游客访问（仅支持低分辨率）"}
	fmt.Println("🔅 当前未登录，请选择是否登录: ")
	for index, item := range items {
		fmt.Printf("  %d. %s\n", index+1, item)
	}
	fmt.Printf("> 请输入操作序号 [%d-%d]: ", 1, len(items))
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
