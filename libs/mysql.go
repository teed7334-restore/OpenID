package libs

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

//MySQL 相關設定
type MySQL struct {
	Db *gorm.DB
}

//New 建構式
func (my MySQL) New() *MySQL {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	charset := os.Getenv("MYSQL_CHARSET")
	parseTime := os.Getenv("MYSQL_PARSETIME")
	loc := os.Getenv("MYSQL_LOC")
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=%s&loc=%s", user, password, host, database, charset, parseTime, loc)
	var err error
	my.Db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return &my
}
