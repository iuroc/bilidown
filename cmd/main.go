package main

import (
	"fmt"
	"log"

	"github.com/iuroc/bilidown"
	"github.com/manifoldco/promptui"
)

func main() {
	cookieValue := promptLogin()
	fmt.Println(cookieValue)
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
				break
			}
		}
	}
	return cookieValue
}

func shouldLogin() bool {
	items := []string{"登录账号（支持全部分辨率）", "游客访问（仅支持低分辨率）"}
	prompt := promptui.Select{
		Items:     items,
		HideHelp:  true,
		Templates: &promptui.SelectTemplates{Label: "当前未登录，请选择操作（方向键选择）"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalln(err)
	}
	return result == items[0]
}
