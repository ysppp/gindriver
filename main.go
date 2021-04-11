package main

import (
	"fmt"
	"gindriver/config"
	"gindriver/models"
	"gindriver/router"
	"gindriver/utils"
	"github.com/jinzhu/configor"
	"time"
)

func main() {
	var err error
	err = configor.
		New(&configor.Config{
			AutoReload:         true,
			AutoReloadInterval: time.Second,
			AutoReloadCallback: func(config interface{}) {
				fmt.Println("config reloaded")
				err = utils.InitDatabase()
				if err != nil {
					fmt.Printf("err: %s\n", err)
				}
				if utils.Database != nil && err == nil {
					utils.Database.AutoMigrate(&models.User{})
				}
			}}).
		Load(&config.Config, "config/config.yml")

	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	err = utils.InitDatabase()
	if err != nil {
		fmt.Printf("[InitDatabase] err: %s\n", err)
	}

	sqlDb, err := utils.Database.DB()
	if err != nil {
		fmt.Printf("[DB] err: %s\n", err)
	}

	defer sqlDb.Close()

	if utils.Database != nil && err == nil {
		utils.Database.AutoMigrate(&models.User{})
		utils.Database.AutoMigrate(&models.FileFolder{})
		utils.Database.AutoMigrate(&models.File{})
		utils.Database.AutoMigrate(&models.Share{})
		utils.Database.AutoMigrate(&models.FileStore{})
	}

	err = utils.InitWebAuthn()
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Printf("%s: running on addr: %s\n", config.Config.AppName, config.Config.ListenAddr)
	app := router.InitRouter()

	err = app.Run(config.Config.ListenAddr)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
}
