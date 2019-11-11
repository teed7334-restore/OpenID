package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/teed7334-restore/oauth_server/libs"
)

func main() {
	doOauth()
}

func doOauth() {
	o := libs.OAuth{}.New()
	o.APIs()
}
