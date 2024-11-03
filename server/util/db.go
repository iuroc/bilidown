package util

import "database/sql"

func CreateLog(db *sql.DB, content string) error {
	_, err := db.Exec(`INSERT INTO "log" ("content") VALUES (?)`, content)
	return err
}
