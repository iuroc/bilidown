package bilidown

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Login 调用浏览器登录并返回 SESSDATA
func Login() (*network.Cookie, error) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.Flag("headless", false),
	)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.bilibili.com/"),
	)
	if err != nil {
		return nil, err
	}
	var loginCookie *network.Cookie
	for {
		time.Sleep(time.Second)
		err = chromedp.Run(ctx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				cookies, err := network.GetCookies().Do(ctx)
				for _, cookie := range cookies {
					if cookie.Name == "SESSDATA" {
						loginCookie = cookie
						return nil
					}
				}
				return err
			}),
		)
		if err != nil {
			return nil, err
		}
		if loginCookie != nil {
			break
		}
	}
	return loginCookie, nil
}

// SaveCookie 将 Cookie 以 JSON 格式保存到 cookie 文件中
func SaveCookie(cookie *network.Cookie, cookieSavePath string) {
	result, err := json.Marshal(cookie)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(cookieSavePath, result, 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

// GetCookieValue 获取文件中保存的可用 Cookie
func GetCookieValue(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return "", err
	} else if err != nil {
		log.Fatalln(err)
	}
	var cookie network.Cookie
	err = json.Unmarshal(data, &cookie)
	if err != nil {
		return "", errors.New("cookie 文件内容格式错误")
	}
	if cookie.Name == "SESSDATA" && expiresToTime(cookie.Expires).After(time.Now()) {
		return cookie.Value, nil
	} else {
		return "", errors.New("无可用 Cookie 或 Cookie 过期")
	}
}

// ExpiresToTime 将 network.Cookie.Expires 转换为 Time
func expiresToTime(expires float64) time.Time {
	seconds := int64(expires)
	nanos := int64((expires - float64(seconds)) * 1e9)
	return time.Unix(seconds, nanos)
}
