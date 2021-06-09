package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//声明数据库连接全局变量
var (
	dbConn *sql.DB
	err    error
)

//init方法在包被调用的时候第一个被执行
func init() {
	dbConn, err = sql.Open("mysql", "root:120120@tcp(localhost:3306)/video?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
}
