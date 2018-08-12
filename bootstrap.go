package main

import (
	"os"

	"github.com/ketan/gocd-golang-bootstrapper/agent"
	"github.com/ketan/gocd-golang-bootstrapper/config"
	"github.com/ketan/gocd-golang-bootstrapper/log"
)

const workingDirectory = "./go"

func main() {
	log.Init()

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
