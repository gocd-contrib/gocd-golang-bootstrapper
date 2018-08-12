package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	goEAServerURLEnv                = "GO_EA_SERVER_URL"
	goEAAutoRegisterKeyEnv          = "GO_EA_AUTO_REGISTER_KEY"
	goEAAutoRegisterEnvironmentsEnv = "GO_EA_AUTO_REGISTER_ENVIRONMENT"
	goEAAutoRegisterAgentIDEnv      = "GO_EA_AUTO_REGISTER_ELASTIC_AGENT_ID"
	goEAAutoRegisterPluginIDEnv     = "GO_EA_AUTO_REGISTER_ELASTIC_PLUGIN_ID"
	goEASSLRootCertFileEnv          = "GO_EA_SSL_ROOT_CERT_FILE"
	goEAGUIDEnv                     = "GO_EA_GUID"
	goEASSLNoVerifyEnv              = "GO_EA_SSL_NO_VERIFY"
)

// GoServerURL reads the env GO_EA_SERVER_URL
func GoServerURL() (string, error) {
	v := os.Getenv(goEAServerURLEnv)

	if strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("variable %s must be set, and should point to the URL of the go server. Example GO_EA_SERVER_URL=https://192.168.0.100:8154/go", goEAServerURLEnv)
	}

	return v, nil
}

// GoEAAutoRegisterContents returns the string to be written in autoregister.properties file
// of the elastic agent
func GoEAAutoRegisterContents() (string, error) {
	autoRegisterKey, err := getEnvAutoregisterEnvAndAssertNotEmpty(goEAAutoRegisterKeyEnv)
	if err != nil {
		return "", err
	}
	agentID, err := getEnvAutoregisterEnvAndAssertNotEmpty(goEAAutoRegisterAgentIDEnv)
	if err != nil {
		return "", err
	}
	pluginID, err := getEnvAutoregisterEnvAndAssertNotEmpty(goEAAutoRegisterPluginIDEnv)
	if err != nil {
		return "", err
	}
	autoRegisterContents := fmt.Sprintf(`
	agent.auto.register.key=%s
	agent.auto.register.environments=%s
	agent.auto.register.elasticAgent.agentId=%s
	agent.auto.register.elasticAgent.pluginId=%s
	`, autoRegisterKey,
		os.Getenv(goEAAutoRegisterEnvironmentsEnv),
		agentID,
		pluginID)

	return autoRegisterContents, nil
}

// InsecureSkipVerify reads GO_EA_SSL_NO_VERIFY flag
func InsecureSkipVerify() bool {
	sslVerify := os.Getenv(goEASSLNoVerifyEnv)

	return (strings.TrimSpace(sslVerify) == "false")
}

// HasSpecifiedRootCAs checks if any Root CA file has been certified
func HasSpecifiedRootCAs() bool {
	return !(strings.TrimSpace(RootCertFile()) == "")
}

// RootCertFile reads GO_EA_SSL_ROOT_CERT_FILE env value
func RootCertFile() string {
	return os.Getenv(goEASSLRootCertFileEnv)
}

// GoEAUUID reads GO_EA_GUID env value
func GoEAUUID() string {
	return os.Getenv(goEAGUIDEnv)
}

func getEnvAutoregisterEnvAndAssertNotEmpty(envName string) (string, error) {
	v := os.Getenv(envName)
	if strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("variable '%s' must be set. See https://docs.go.cd/current/advanced_usage/agent_auto_register.html for more information", envName)
	}
	return v, nil
}
