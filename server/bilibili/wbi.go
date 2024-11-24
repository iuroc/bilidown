package bilibili

import "bilidown/util"

func (client *BiliClient) getWebKey() (imgKey string, subKey string, err error) {
	db := util.MustGetDB()
	defer db.Close()
	imgKey, subKey, err = util.GetWbiKey(db)
	if err != nil {
		return "", "", err
	}
	if imgKey == "" || subKey == "" {
		// 更新数据库中的 WebKey
		imgKey, subKey, err = client.GetWbiKey()
		if err = util.SaveWbiKey(db, imgKey, subKey); err != nil {
			return "", "", err
		}
	}
	return imgKey, subKey, nil
}
