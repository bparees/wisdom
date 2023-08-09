# OpenShift Wisdom Project

## Overview
Experiments in creating a service layer to sit between UX clients and backend LLMs.  

Intended to provide capabilities such as:
  * prompt engineering/augumentation (adding additional context to the user supplied prompt, rejecting poor quality prompts)
  * response validating/linting/filtering (stripping sensitive information, yaml validation)
  * allowing selection from multiple backend models
  * response attribution (match response to training data and cite the reference)
  * record user feedback about response quality

## Building and Usage

### Build the binary
$ cd cmd/wisdom
$ go build .

### Define a config file
See https://github.com/bparees/wisdom/blob/main/config/sample.cfg.yaml for an example.

### Run a server
$ ./wisdom serve --config path/to/config.yaml

### Do a single inference
$ ./wisdom infer --config path/to/config.yaml --prompt "write a deployment yaml for the registry.redhat.io/rhel9/redis-6:latest image with 3 replicas"