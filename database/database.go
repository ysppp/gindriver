package database

import (
	"fmt"
	"gindriver/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

var Database *gorm.DB

func InitDatabase() (err error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?%s",
		config.Config.DB.User,
		config.Config.DB.Pass,
		config.Config.DB.Host,
		strconv.Itoa(int(config.Config.DB.Port)),
		config.Config.DB.Name,
		config.Config.DB.Param)
	fmt.Printf("DSN: %s", dsn)
	Database, err = gorm.Open("mysql", dsn)

	return err
}
