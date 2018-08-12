package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/ketan/gocd-golang-bootstrapper/env"
	"github.com/ketan/gocd-golang-bootstrapper/log"
)

type agentChecksums struct {
	agentMd5    string
	pluginsMd5  string
	launcherMd5 string
	tfsImplMd5  string
}

// Start setups up the prerequisites and runs the agent jar
func Start() error {
	url, err := env.GoServerURL()
	if err != nil {
		return err
	}
	checksumURL := fmt.Sprintf("%s/admin/latest-agent.status", url)
	for {
		log.Debugf("Getting checksums from %s", checksumURL)
		chk, err := getChecksums(checksumURL)
		if err != nil {
			return err
		}
		log.Debugf("agent.jar                     - %s", chk.agentMd5)
		log.Debugf("agent-plugins.zip             - %s", chk.pluginsMd5)
		log.Debugf("tfs-impl.jar                  - %s", chk.tfsImplMd5)
		log.Debugf("agent-launcher.jar (not used) - %s", chk.launcherMd5)

		files := map[string]string{
			fmt.Sprintf("%s/admin/agent", url):             "agent.jar",
			fmt.Sprintf("%s/admin/agent-plugins.zip", url): "agent-plugins.zip",
			fmt.Sprintf("%s/admin/tfs-impl.jar", url):      "tfs-impl.jar",
		}
		err = downloadFiles(files)
		if err != nil {
			return err
		}
		err = startAgent(url, chk)
		if err != nil {
			return err
		}
	}
}

func startAgent(goServerURL string, checksums *agentChecksums) error {
	cmd := exec.Command("java",
		"-Dcruise.console.publish.interval=10",
		"-Xms128m",
		"-Xmx256m",
		"-Djava.security.egd=file:/dev/./urandom",
		fmt.Sprintf("-Dagent.plugins.md5=%s", checksums.pluginsMd5),
		fmt.Sprintf("-Dagent.binary.md5=%s", checksums.agentMd5),
		fmt.Sprintf("-Dagent.launcher.md5=%s", checksums.launcherMd5),
		fmt.Sprintf("-Dagent.tfs.md5=%s", checksums.tfsImplMd5),
		"-jar",
		"agent.jar",
		"-serverUrl",
		goServerURL,
		"-sslVerificationMode",
	)

	if env.InsecureSkipVerify() {
		cmd.Args = append(cmd.Args, "NONE")
	} else {
		cmd.Args = append(cmd.Args, "FULL")

		if env.HasSpecifiedRootCAs() {
			cmd.Args = append(cmd.Args, "-rootCertFile", env.RootCertFile())
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// we remove the environment variables that are needed by us, but don't need to be passed onto the agent process
	re := regexp.MustCompile("^(GO_EA_).*")
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
		return fmt.Errorf("Could not launch agent process %s", err.Error())
	}

	log.Infof("Agent process launched PID(%d), waiting for it to exit.", cmd.Process.Pid)
	err = cmd.Wait()
	log.Infof("Agent process PID(%d) exited with %v.\nSleeping 10 seconds, before starting over again", cmd.Process.Pid, cmd.ProcessState)
	time.Sleep(10 * time.Second)
	return nil
}

func downloadFiles(f map[string]string) error {
	r, err := rootCAs()
	if err != nil {
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: env.InsecureSkipVerify(),
			RootCAs:            r,
		},
	}

	client := &http.Client{
		Transport: tr,
	}
	for dest, url := range f {
		log.Infof("Downloading file %s from %s", dest, url)
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("Could not fetch URL %s. %v", url, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("the URL %s returned %d", url, resp.StatusCode)
		}

		out, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("Unable to write to file %s. %s", dest, err.Error())
		}
		defer out.Close()
		io.Copy(out, resp.Body)
		log.Infof("Finished downloading file %s from %s", dest, url)
	}
	return nil
}

func getChecksums(url string) (*agentChecksums, error) {
	r, err := rootCAs()
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: env.InsecureSkipVerify(),
			RootCAs:            r,
		},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Head(url)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch URL %s. %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the URL %s returned %d", url, resp.StatusCode)
	}
	return &agentChecksums{agentMd5: resp.Header.Get("Agent-Content-MD5"),
		pluginsMd5:  resp.Header.Get("Agent-Plugins-Content-MD5"),
		launcherMd5: resp.Header.Get("Agent-Launcher-Content-MD5"),
		tfsImplMd5:  resp.Header.Get("TFS-SDK-Content-MD5"),
	}, nil
}

func rootCAs() (*x509.CertPool, error) {
	if !env.HasSpecifiedRootCAs() {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("Couldn't load system certificate pool, %s", err.Error())
		}
		return certPool, nil
	}
	cert, err := ioutil.ReadFile(env.RootCertFile())
	if err != nil {
		return nil, fmt.Errorf("Couldn't load file, %s", err.Error())
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	return certPool, nil
}
