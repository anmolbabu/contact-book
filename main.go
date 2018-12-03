package main

import (
	"fmt"
	"os"

	_ "github.com/heroku/x/hmetrics/onload"

	"github.com/anmolbabu/contact-book/api"
	"github.com/anmolbabu/contact-book/config"
	"github.com/anmolbabu/contact-book/dao"
)

func main() {

	config := config.GetConfig()

	db, err := dao.Init(config)
	if err != nil {
		fmt.Printf("Failed to initialise db with config %+v. Error: %v\n", config, err)
		os.Exit(1)
	}
	defer db.Cleanup()

	api.NewRouter(config)

}
