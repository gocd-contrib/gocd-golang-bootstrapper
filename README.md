# Native Binary GoCD agent bootstrapper

This is primarily intended to be used by elastic agents to reduce startup time of agents and reduce the memory footprint on the agent machine, since this runs as a native executable.

# Usage

## Download

Get your OS/Architecture specific binary from `https://github.com/gocd-contrib/gocd-golang-bootstrapper/releases/latest`

## Run

```shell
./go-bootstrapper
```

# Supported environment variables

The program uses environment variables to configure the agent to make it easy to embed with docker.

## General options:

| Environment              | Required | Description                                                                                                              |
| :----------------------- | :------- | :----------------------------------------------------------------------------------------------------------------------- |
| `GO_EA_GUID`             | No       | The contents of `config/guid.txt` file that contains the agent identifier.                                               |
| `GO_EA_SERVER_URL`       | Yes      | The base url to the GoCD server.                                                                                         |
| `GO_EA_DUMP_ENVIRONMENT` | No       | Whether environment variables should be dumped to the log. Turned off by default for security.                           |
| `GO_EA_JVM_ARGS`         | No       | JSON formatted list of args that should be passed as JVM args to the agent process. Example: `["-Dfoo=bar", "-Xmx256m"]` |
| `GO_EA_ROOT_DIR`         | No       | The directory where the gocd agent should execute. Defaults to `/go`.                                                    |

## Autoregistration options:

| Environment                             | Required | Description                                                                  |
| :-------------------------------------- | :------- | :--------------------------------------------------------------------------- |
| `GO_EA_AUTO_REGISTER_KEY`               | Yes      | The GoCD agent auto register key.                                            |
| `GO_EA_AUTO_REGISTER_ENVIRONMENT`       | No       | The name of the environment that the agent should autoregister with.         |
| `GO_EA_AUTO_REGISTER_ELASTIC_AGENT_ID`  | Yes      | The elastic agent identifier that the agent should autoregister with.        |
| `GO_EA_AUTO_REGISTER_ELASTIC_PLUGIN_ID` | Yes      | The elastic agent plugin identifier that the agent should autoregister with. |

## SSL options:

| Environment                | Required | Description                                                                                                  |
| :------------------------- | :------- | :----------------------------------------------------------------------------------------------------------- |
| `GO_EA_SSL_NO_VERIFY`      | No       | Whether ssl verification should be turned off. Defaults to `false`.                                          |
| `GO_EA_SSL_ROOT_CERT_FILE` | No       | The path to the file containing root CA certificates. Defaults to the file provided by the operating system. |

## Building instructions

### Using `go get`

```
go install github.com/gocd-contrib/gocd-golang-bootstrapper@latest
```

### Using `make`

```
$ make all
```
