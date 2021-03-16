package utils

import (
	"fmt"
	"gindriver/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

var Database *gorm.DB

func InitDatabase() (err error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		config.Config.DB.User,
		config.Config.DB.Pass,
		config.Config.DB.Host,
		strconv.Itoa(int(config.Config.DB.Port)),
		config.Config.DB.Name)
	//config.Config.DB.Param)
	fmt.Printf("DSN: %s\n", dsn)
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}
