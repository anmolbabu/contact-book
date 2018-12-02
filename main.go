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

	_, err := dao.Init(config)
	if err != nil {
		fmt.Printf("Failed to initialise db with config %+v. Error: %v\n", config, err)
		os.Exit(1)
	}

	api.NewRouter(config)

}
