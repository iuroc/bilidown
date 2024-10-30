package task

import (
	"database/sql"
	"os"
)

type TaskOption struct {
	Bvid   string `json:"bvid"`
	Cid    int    `json:"cid"`
	Format int    `json:"format"`
}

func CreateTask(db *sql.DB, option TaskOption, filepath string) error {
	_, err := db.Exec(`INSERT INTO "task" ("bvid", "cid", "format", "filepath") VALUES (?, ?, ?, ?)`, option.Bvid, option.Cid, option.Format, filepath)
	return err
}

func GetCurrentFolder(db *sql.DB) (string, error) {
	var filepath string
	err := db.QueryRow(`SELECT "value" FROM "field" WHERE "name" = "download_folder"`).Scan(&filepath)
	if err != nil {
		return "", err
	}
	return filepath, nil
}

// SaveDownloadFolder 保存下载路径，不存在则自动创建
func SaveDownloadFolder(db *sql.DB, downloadFolder string) error {
	_, err := os.Stat(downloadFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(downloadFolder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return err
	}
	_, err = db.Exec(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES ("download_folder", ?)`, downloadFolder)
	return err
}
