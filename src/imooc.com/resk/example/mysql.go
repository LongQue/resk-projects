package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	dsName := "root:123456@tcp(127.0.0.1:3306)/resk?charset=utf8&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	//最大保持连接数
	db.SetMaxIdleConns(2)
	//最大连接数
	db.SetMaxOpenConns(3)
	//默认8小时，不可大于
	db.SetConnMaxLifetime(7 * time.Hour)

	fmt.Println(db.Query("select now()"))
}
