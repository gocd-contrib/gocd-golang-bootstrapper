package main

import (
	"crypto/tls"
	"fmt"
	"github.com/op/go-logging"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var log = logging.MustGetLogger("gocd-golang-bootstrapper")
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05.999Z-07:00} [Bootstrap] [%{level:.4s}]%{color:reset} %{message}`,
)

func main() {
	var err error = nil
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(backendFormatter)

	autoRegisterContents, goServerUrl := os.Getenv("AUTO_REGISTER_CONTENTS"), os.Getenv("GO_SERVER_URL")

	if strings.TrimSpace(autoRegisterContents) == "" {
		log.Critical(`The variable AUTO_REGISTER_CONTENTS must be set, and should contain the contents of the autoregister.properties file. See https://docs.go.cd/current/advanced_usage/agent_auto_register.html for more information.`)
		os.Exit(1)
	}

	if strings.TrimSpace(goServerUrl) == "" {
		log.Critical("The variable GO_SERVER_URL must be set, and should point to the URL of the go server. Example GO_SERVER_URL=https://192.168.0.100:8154/go")
		os.Exit(1)
	}

	err = os.Chdir("/go")
	if err != nil {
		log.Criticalf("Could not change working directory to /go. %v", err)
		os.Exit(1)
	}

	err = os.RemoveAll("config")
	if err != nil {
		log.Criticalf("There was an error removing the 'config' directory. %v", err)
		os.Exit(1)
	}

	err = os.Mkdir("config", 0750)
	if err != nil {
		log.Criticalf("There was an error creating the 'config' directory. %v", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile("config/autoregister.properties", []byte(autoRegisterContents), 0600)

	logContents := os.Getenv("LOG_FILE_CONTENTS")
	if strings.TrimSpace(logContents) == "" {
		hostName, err := os.Hostname()
		if err != nil {
			log.Warning("Could not detect hostname, assuming HOSTNAME environment. %v", err)
			hostName = os.Getenv("HOSTNAME")
			err = nil
		}

		logContents = fmt.Sprintf(`
# default to INFO logging on stdout and tcp
log4j.rootLogger=INFO, stdout, tcp

# write logs to stdout
log4j.appender.stdout=org.apache.log4j.ConsoleAppender
log4j.appender.stdout.layout=org.apache.log4j.PatternLayout
log4j.appender.stdout.layout.conversionPattern=%%d{ISO8601} [%%-9t] %%-5p %%-16c{4}:%%L %%x- %%m%%n

# write logs to log server
log4j.appender.tcp=org.apache.log4j.net.SocketAppender
log4j.appender.tcp.RemoteHost=%s
log4j.appender.tcp.ReconnectionDelay=10000
log4j.appender.tcp.Application=%s
`, os.Getenv("LOGS_HOST"), hostName)
	}

	err = ioutil.WriteFile("log4j.properties", []byte(logContents), 0600)
	if err != nil {
		log.Warning("Could not write log4j.properites, continuing.")
		err = nil
	}
	err = ioutil.WriteFile("go-agent-log4j.properties", []byte(logContents), 0600)
	if err != nil {
		log.Warning("Could not write go-agent-log4j.properties, continuing.")
		err = nil
	}

	checksumUrl := fmt.Sprintf("%s/admin/latest-agent.status", goServerUrl)
	log.Debugf("Getting checksums from %s", checksumUrl)
	agentMd5, agentPluginsMd5, agentLauncherMd5 := getChecksums(checksumUrl)
	log.Debugf("agent.jar                     - %s", agentMd5)
	log.Debugf("agent-plugins.zip             - %s", agentPluginsMd5)
	log.Debugf("agent-launcher.jar (not used) - %s", agentLauncherMd5)

	for {
		downloadFile(fmt.Sprintf("%s/admin/agent", goServerUrl), "agent.jar")
		downloadFile(fmt.Sprintf("%s/admin/agent-plugins.zip", goServerUrl), "agent-plugins.zip")

		startAgent(goServerUrl, agentMd5, agentPluginsMd5, agentLauncherMd5)
	}
}

func startAgent(goServerUrl string, agentMd5 string, agentPluginsMd5 string, agentLauncherMd5 string) {
	cmd := exec.Command("java",
		"-Dcruise.console.publish.interval=10",
		"-Xms128m",
		"-Xmx256m",
		"-Djava.security.egd=file:/dev/./urandom",
		fmt.Sprintf("-Dagent.plugins.md5=%s", agentPluginsMd5),
		fmt.Sprintf("-Dagent.binary.md5=%s", agentMd5),
		fmt.Sprintf("-Dagent.launcher.md5=%s", agentLauncherMd5),
		"-jar",
		"agent.jar",
		"-serverUrl",
		goServerUrl,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// we remove the environment variables that are needed by us, but don't need to be passed onto the agent process
	re := regexp.MustCompile("^(AUTO_REGISTER_CONTENTS|LOGS_HOST)=.*")
	filteredEnv := make([]string, 0)
	for _, elem := range os.Environ() {
		if !re.MatchString(elem) {
			filteredEnv = append(filteredEnv, elem)
		}
	}

	cmd.Env = filteredEnv

	log.Infof("Launching command %v with environment %v", cmd.Args, cmd.Env)
	err := cmd.Start()

	if err != nil {
		log.Critical("Could not launch agent process %v", err)
		os.Exit(1)
	}

	log.Infof("Agent process launched PID(%d), waiting for it to exit.", cmd.Process.Pid)
	err = cmd.Wait()
	log.Infof("Agent process PID(%d) exited with %v.", cmd.Process.Pid, cmd.ProcessState)

	log.Info("Sleeping 10 seconds, before starting over again")
	time.Sleep(10 * time.Second)
}

func downloadFile(url string, dest string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	checkHttpResponse(url, resp, err)
	defer resp.Body.Close()
	out, err := os.Create(dest)

	if err != nil {
		log.Critical("Unable to write to file %s. %v", dest, err)
		os.Exit(1)
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}

func getChecksums(url string) (string, string, string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Head(url)
	checkHttpResponse(url, resp, err)
	resp.Body.Close()
	return resp.Header.Get("Agent-Content-MD5"), resp.Header.Get("Agent-Plugins-Content-MD5"), resp.Header.Get("Agent-Launcher-Content-MD5")
}

func checkHttpResponse(url string, resp *http.Response, err error) {
	if err != nil {
		log.Criticalf("Could not fetch URL %s. %v", url, err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		log.Criticalf("The URL %s returned %d.", url, resp.StatusCode)
		os.Exit(1)
	}
}
