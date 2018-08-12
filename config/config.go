package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ketan/gocd-golang-bootstrapper/env"
)

const logbackContents = `
<?xml version="1.0" encoding="UTF-8"?>
<included>
	  <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
    <!-- encoders are assigned the type
         ch.qos.logback.classic.encoder.PatternLayoutEncoder by default -->
    <encoder>
      <pattern>
      	${gocd.agent.logback.defaultPattern:-%%date{ISO8601} %%-5level [%%thread] %%logger{0}:%%line - %%msg%%n}
      </pattern>
    </encoder>
  </appender>

  <root>
    <appender-ref ref="STDOUT" />
  </root>
</included>
`

// Write writes setup config in the agent: logback-config.xml, agent UUID and autoregister.properties
func Write() error {
	err := clean()
	if err != nil {
		return err
	}

	err = writeLogbackFileContents()
	if err != nil {
		return err
	}

	err = writeUUIDContents()
	if err != nil {
		return err
	}

	err = writeAutoregisterContents()
	if err != nil {
		return err
	}

	return nil
}

func clean() error {
	err := os.RemoveAll("config")
	if err != nil {
		return fmt.Errorf("There was an error removing the 'config' directory. %s", err.Error())
	}

	err = os.Mkdir("config", 0750)
	if err != nil {
		return fmt.Errorf("There was an error creating the 'config' directory. %s", err.Error())
	}
	return nil
}

func writeLogbackFileContents() error {
	err := ioutil.WriteFile("config/go-agent-logback-include.xml", []byte(strings.TrimSpace(logbackContents)), 0600)
	if err != nil {
		return fmt.Errorf("Unable to write config/go-agent-logback-include.xml, %s", err.Error())
	}
	return nil
}

func writeUUIDContents() error {
	uuidContents := env.GoEAUUID()

	if strings.TrimSpace(uuidContents) == "" {
		return nil
	}

	err := ioutil.WriteFile("config/guid.txt", []byte(uuidContents), 0600)
	if err != nil {
		return fmt.Errorf("Could not write config/guid.txt, %s", err.Error())
	}
	return nil
}

func writeAutoregisterContents() error {
	v, err := env.GoEAAutoRegisterContents()
	if err != nil {
		return fmt.Errorf("Could not write config/autoregister.properties, %s", err.Error())
	}

	err = ioutil.WriteFile("config/autoregister.properties", []byte(v), 0600)
	if err != nil {
		return fmt.Errorf("Could not write config/autoregister.properties, %s", err.Error())
	}

	return nil
}
