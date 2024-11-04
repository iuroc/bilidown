package util

import (
	"database/sql"
	"strings"
)

func CreateLog(db *sql.DB, content string) error {
	_, err := db.Exec(`INSERT INTO "log" ("content") VALUES (?)`, content)
	return err
}

func GetFields(db *sql.DB, names ...string) (map[string]string, error) {
	if len(names) == 0 {
		return nil, nil
	}
	row, err := db.Query(`SELECT "name", "value" FROM "field" WHERE "name" IN (?)`, strings.Join(names, ", "))
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
