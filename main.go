package main

import (
	config2 "github.com/rodrinoblega/microblogging/config"
	"github.com/rodrinoblega/microblogging/setup"
	"log"
	"os"
)

func main() {

	envConf := config2.Load(os.Getenv("ENV"))

	var appDependencies *setup.AppDependencies
	appDependencies = setup.InitializeDependencies(envConf)

	log.Printf("running %s environment", envConf.Env)
	router := SetupRouter(appDependencies)

	router.Run(":8080")
}
