package bilibili

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/iuroc/bilidown/server/util"
)

type BiliClient struct {
	SESSDATA string
}

// SimpleGET 简单的 GET 请求
func (client *BiliClient) SimpleGET(_url string, params map[string]string) (*http.Response, error) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	_client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(nil),
		},
	}
	request, err := http.NewRequest("GET", _url+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	request.Header = client.MakeHeader()
	return _client.Do(request)
}

// MakeHeader 生成请求头
func (client *BiliClient) MakeHeader() http.Header {
	header := http.Header{}
	header.Set("Cookie", "SESSDATA="+client.SESSDATA)
	header.Set("User-Agent", "Mozilla/5.0")
	header.Set("Referer", "https://www.bilibili.com")
	return header
}

// CheckLogin 检查是否已经登录
func (client *BiliClient) CheckLogin() (bool, error) {
	response, err := client.SimpleGET("https://api.bilibili.com/x/space/myinfo", nil)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	body := BaseResV2{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return false, err
	}
	if body.Code != 0 {
		return false, errors.New(body.Message)
	}
	return body.Success(), nil
}

// NewQRInfo 获取登录二维码信息
func (client *BiliClient) NewQRInfo() (*QRInfo, error) {
	response, err := client.SimpleGET("https://passport.bilibili.com/x/passport-login/web/qrcode/generate", nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body := BaseResV2{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, errors.New(body.Message)
	}
	qrInfo := QRInfo{}
	err = json.Unmarshal(body.Data, &qrInfo)
	if err != nil {
		return nil, err
	}
	return &qrInfo, nil
}

// GetQRStatus 获取二维码状态
func (client *BiliClient) GetQRStatus(qrKey string) (qrStatus *QRStatus, sessdata string, err error) {
	params := map[string]string{
		"qrcode_key": qrKey,
	}
	response, err := client.SimpleGET("https://passport.bilibili.com/x/passport-login/web/qrcode/poll", params)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()
	body := BaseResV2{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, "", err
	}
	if body.Code != 0 {
		return nil, "", errors.New(body.Message)
	}
	qrStatus = &QRStatus{}
	err = json.Unmarshal(body.Data, &qrStatus)
	if err != nil {
		return nil, "", err
	}
	if qrStatus.Code != 0 {
		return qrStatus, "", nil
	}
	sessdata, err = GetCookieValue(response.Cookies(), "SESSDATA")
	if err != nil {
		return nil, "", err
	}
	return qrStatus, sessdata, nil
}

// GetCookieValue 获取指定 Name 的 Cookie 值
func GetCookieValue(cookies []*http.Cookie, name string) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value, nil
		}
	}
	return "", errors.New("cookie with name " + name + " not found")
}

// SaveSessdata 保存 SESSDATA
func SaveSessdata(db *sql.DB, sessdata string) error {
	util.SqliteLock.Lock()
	_, err := db.Exec(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES ("sessdata", ?)`, sessdata)
	util.SqliteLock.Unlock()
	return err
}

// GetSessdata 获取 SESSDATA
func GetSessdata(db *sql.DB) (string, error) {
	util.SqliteLock.Lock()
	row := db.QueryRow(`SELECT "value" FROM "field" WHERE "name" = "sessdata"`)
	util.SqliteLock.Unlock()
	var sessdata string
	err := row.Scan(&sessdata)
	return sessdata, err
}
