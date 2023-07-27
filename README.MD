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

$ go build .

# run a server
$ ./wisdom --email <your email address> --apikey <your IBM openshift LLM workbench apikey>

# do a single inference
$ ./wisdom infer --email <email address> --apikey <apikey> --prompt "write a deployment yaml for the registry.redhat.io/rhel9/redis-6:latest image with 3 replicas"