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
			AutoReloadInterval: time.Minute}).
		Load(&config.Config, "config/config.yml")

	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	err = utils.InitDatabase()
	defer utils.Database.Close()
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	if utils.Database != nil {
		utils.Database.AutoMigrate(&models.User{})
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
