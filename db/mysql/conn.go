package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	storeconst "github.com/luoruofeng/private-store/const"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", storeconst.DB_ADDR)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("failed to connect to mysql err:%s\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("mysql ping success!")
	}
}

//DBConn : return connected db
func DBConn() *sql.DB {
	return db
}
