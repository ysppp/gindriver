package main

import (
	"fmt"
	"gindriver/config"
	"gindriver/database"
	"gindriver/router"
	"github.com/jinzhu/configor"
)

func main() {
	var err error
	err = configor.Load(&config.Config, "config/config.yml")
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	err = database.InitDatabase()
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	fmt.Printf("%s: running on addr: %s\n", config.Config.AppName, config.Config.ListenAddr)
	app := router.InitRouter()
	err = app.Run(config.Config.ListenAddr)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
}
