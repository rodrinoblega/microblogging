package main

import (
	"github.com/rodrinoblega/microblogging/config"
	"github.com/rodrinoblega/microblogging/setup"
	"log"
	"os"
)

func main() {

	envConf := config.Load(os.Getenv("ENV"))

	var appDependencies *setup.AppDependencies
	appDependencies = setup.InitializeDependencies(envConf)

	log.Printf("running %s environment", envConf.Env)
	router := SetupRouter(appDependencies)

	router.Run(":8080")
}
