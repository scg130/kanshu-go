package repo

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var x *xorm.Engine

func init() {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWD")
	db := os.Getenv("MYSQL_NOVEL_DB")
	mysql_log := os.Getenv("MYSQL_LOG") == "true"
	port := os.Getenv("MYSQL_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", user, pass, host, port, db)
	var err error
	x, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("mysql connect err:%v", err))
	}
	if err := x.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("mysql ping err:%v", err))
	}
	x.ShowSQL(mysql_log)

	x.SetMaxIdleConns(5)

	x.SetMaxOpenConns(5)

}
