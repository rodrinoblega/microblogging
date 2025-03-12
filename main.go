package main

import (
	config2 "github.com/rodrinoblega/microblogging/config"
	"github.com/rodrinoblega/microblogging/setup"
	"os"
)

func main() {

	envConf := config2.Load(os.Getenv("ENV"))

	var appDependencies *setup.AppDependencies
	switch envConf.Env {
	case "local":
		appDependencies = setup.InitializeLocalDependencies(envConf)
	}

	router := SetupRouter(appDependencies)

	router.Run(":8080")
}
