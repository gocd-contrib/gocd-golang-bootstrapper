package main

import (
	"os"

	"github.com/gocd-contrib/gocd-golang-bootstrapper/agent"
	"github.com/gocd-contrib/gocd-golang-bootstrapper/config"
	"github.com/gocd-contrib/gocd-golang-bootstrapper/env"
	"github.com/gocd-contrib/gocd-golang-bootstrapper/log"
)

func main() {
	workingDirectory := env.GoRootDir()
	log.Init()

	log.Infof("Executing bootstrapper in %s", workingDirectory)

	err := os.Chdir(workingDirectory)
	if err != nil {
		log.Criticalf("Could not change working directory to %s. %s", workingDirectory, err.Error())
		os.Exit(1)
	}

	err = config.Write()
	if err != nil {
		log.Criticalf("There was an error configuring the agent. %s", err.Error())
		os.Exit(1)
	}

	err = agent.Start()
	if err != nil {
		log.Criticalf("There was an error starting the agent. %s", err.Error())
		os.Exit(1)
	}
}
