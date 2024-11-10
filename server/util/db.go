package util

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func CreateLog(db *sql.DB, content string) error {
	SqliteLock.Lock()
	defer SqliteLock.Unlock()
	_, err := db.Exec(`INSERT INTO "log" ("content") VALUES (?)`, content)
	return err
}

func GetFields(db *sql.DB, names ...string) (map[string]string, error) {
	SqliteLock.Lock()
	defer SqliteLock.Unlock()
	if len(names) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(names))
	for i := 0; i < len(names); i++ {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(`SELECT "name", "value" FROM "field" WHERE "name" IN (%s)`, strings.Join(placeholders, ","))

	values := make([]interface{}, len(names))
	for i := 0; i < len(names); i++ {
		values[i] = names[i]
	}

	row, err := db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var name, value string
	fields := make(map[string]string)
	for row.Next() {
		if err := row.Scan(&name, &value); err != nil {
			return nil, err
		}
		fields[name] = value
	}
	return fields, nil
}

func SaveFields(db *sql.DB, data [][2]string) error {
	SqliteLock.Lock()
	defer SqliteLock.Unlock()
	if len(data) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range data {
		_, err = stmt.Exec(d[0], d[1])
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCurrentFolder 获取数据库中的下载保存路径，如果不存在则将默认路径保存到数据库
func GetCurrentFolder(db *sql.DB) (string, error) {
	SqliteLock.Lock()
	defer SqliteLock.Unlock()
	var folder string
	err := db.QueryRow(`SELECT "value" FROM "field" WHERE "name" = 'download_folder'`).Scan(&folder)
	if err != nil && err == sql.ErrNoRows {
		folder, err = GetDefaultDownloadFolder()
		if err != nil {
			return "", err
		}
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return "", err
		}
		err = SaveDownloadFolder(db, folder)
		if err != nil {
			return "", err
		}
		return folder, nil
	}
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}
	return folder, nil
}

// SaveDownloadFolder 保存下载路径，不存在则自动创建
func SaveDownloadFolder(db *sql.DB, downloadFolder string) error {
	SqliteLock.Lock()
	defer SqliteLock.Unlock()
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
	_, err = db.Exec(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES ('download_folder', ?)`, downloadFolder)
	return err
}

var SqliteLock sync.Mutex

func MustGetDB(path ...string) *sql.DB {
	pathStr := ""
	if len(path) == 0 {
		pathStr = "./data.db"
	} else if len(path) > 1 {
		log.Fatalln(errors.New("len(path) <= 1"))
	} else {
		pathStr = path[0]
	}
	db, err := sql.Open("sqlite", pathStr)
	if err != nil {
		log.Fatalln("sql.Open:", err)
	}
	return db
}
