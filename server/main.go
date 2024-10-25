package main

import (
	"bilidown/router"
	"bilidown/util"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	InitTables()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/api/", http.StripPrefix("/api", router.API()))

	fmt.Println("http://127.0.0.1:8098")
	http.ListenAndServe(":8098", nil)
}

// InitTables 初始化数据表
func InitTables() *sql.DB {
	db := util.GetDB()
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "field" (
		"name" TEXT PRIMARY KEY NOT NULL,
		"value" TEXT
	)`)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
