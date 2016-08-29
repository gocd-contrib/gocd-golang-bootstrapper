# GoLang based GoCD agent bootstrapper

This is primarily intended to be used by elastic agents to reduce startup time of agents and reduce the memory footprint on the agent machine.

# !!!WARNING!!!

This binary does not currently support any SSL verification of the server before starting up the agent process. Please use it at your own risk (or submit a PR to fix this :)

## Building instructions

```
$ make all
```
