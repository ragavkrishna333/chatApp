package Dblocal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func LocalDbConnect() (*sql.DB, error) {
	log.Println("LocalDBConnect(+)")
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "ST752", "000ST75228", "192.168.2.5", 3306, "ragavkrishna")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Println("Open Connection failed:", err.Error())
		return db, err
	}
	log.Println("LocalDBConnect(-)")
	return db, nil
}
