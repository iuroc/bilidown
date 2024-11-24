package bilibili

import (
	"bilidown/util"
	"database/sql"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var MIXIN_KEY_ENC_TAB = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11, 36, 20, 34, 44, 52,
}

// GetWebKey 获取最新可用的 WbiKey，自动刷新缓存
func (client *BiliClient) getWbiKey(db *sql.DB) (wbiKey string, err error) {
	wbiKey, err = getWbiKeyFromDB(db)
	if err != nil {
		return "", err
	}
	if wbiKey == "" {
		// 更新数据库中的 WebKey
		wbiKey, err = client.getWbiKeyRemote()
		if err != nil {
			return "", err
		}
		if err = saveWbiKey(db, wbiKey); err != nil {
			return "", err
		}
	}
	return wbiKey, nil
}

func (client *BiliClient) GetMixinKey(db *sql.DB) (string, error) {
	wbiKey, err := client.getWbiKey(db)
	if err != nil {
		return "", err
	}
	var result string
	for _, index := range MIXIN_KEY_ENC_TAB {
		result += string(wbiKey[index])
	}
	return result[:32], nil
}

func WbiSign(params map[string]string, mixinKey string) url.Values {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	wts := strconv.FormatInt(time.Now().Unix(), 10)
	values.Set("wts", wts)
	encodeStr := strings.ReplaceAll(values.Encode(), "+", "%20") + mixinKey
	w_rid := util.MD5Hash(encodeStr)
	values.Set("w_rid", w_rid)
	return values
}

// GetWbiKey 从数据库中获取 wbiKey，如果数据库中不存在记录或记录过期，则返回空字符串但不返回错误
func getWbiKeyFromDB(db *sql.DB) (wbiKey string, err error) {
	fields, err := util.GetFields(db, "wbi_key", "wbi_key_update_at")
	if err != nil {
		return "", err
	}
	// 获取上次更新时间的时间戳，单位是秒
	updateAt, err := strconv.ParseInt(fields["wbi_key_update_at"], 10, 64)
	if err != nil {
		return "", nil
	}
	// 注意 key_update_at 单位是秒，判断上次刷新时间是否超过 1 天
	if time.Now().Unix()-updateAt > 24*60*60 {
		return "", nil
	}
	return fields["wbi_key"], nil
}

// SaveWbiKey 保存 imgKey 和 subKey 到数据库
func saveWbiKey(db *sql.DB, wbiKey string) error {
	return util.SaveFields(db, [][2]string{
		{"wbi_key", wbiKey},
		{"wbi_key_update_at", strconv.FormatInt(time.Now().Unix(), 10)},
	})
}
