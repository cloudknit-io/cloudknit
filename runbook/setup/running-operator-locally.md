# Run Operator locally

## Overview

To speed up local development and improve developer experience, a dev should be able to run/debug operator locally.

## When to use this runbook

When you want to test a piece of code by running the operator locally first

## Prerequisites

1. [Telepresence](https://www.telepresence.io/) - FAST, LOCAL DEVELOPMENT FOR KUBERNETES AND OPENSHIFT MICROSERVICES
2. [IntelliJ](https://www.jetbrains.com/idea/) - OPTIONAL: Integrated Development Environment

## Initial Steps Overview

- [Create a local env file](#create-local-env-file)
- [Proxy your machine to k8s cluster](#proxy-your-machine-to-k8s-cluster)
- [Running the operator](#running-the-operator)
- [Enabling Webhook Service](#enabling-webhook-service)

## Detailed Steps

### Create local env file
1. Run the following script to get the Operator environment variables
	```shell script
	POD_NAME=$(kubectl get pods --namespace zlifecycle-il-operator-system -l "app.kubernetes.io/instance=zlifecycle-il-operator" -o jsonpath="{.items[0].metadata.name}")
	kubectl exec --namespace zlifecycle-il-operator-system -it $POD_NAME -- env
	```
2. Save the environment variables into `PROJECT_ROOT/<environment_name>.env` (ex. `sandbox.env`)
3. Add `DISABLE_WEBHOOKS=true` so it doesn't run the webhook server locally, until we fix the local cert issue

### Proxy your machine to k8s cluster
1. Select kubecontext (ie. sandbox, demo...)
2. Run `telepresence connect`

### Running the operator

#### IntelliJ
1. Install the plugin [EnvFile](https://plugins.jetbrains.com/plugin/7861-envfile)
2. Edit -> Edit configurations -> Add New Configuration -> Go Build -> select `Package` for `Run kind:`
3. Select the `EnvFile` tab -> Enable EnvFile -> Add -> Select your env file for your environment
4. Now you can run/debug your operator code: Run -> Run: '<configuration-name>' | Debug: '<configuration-name>'

#### VS Code

1. Install [delve](https://github.com/go-delve/delve)
1. `mkdir .vscode && touch .vscode/launch.json`
1. Copy and paste the following into `launch.json`
	```
	{
		"version": "0.2.0",
		"configurations": [
			{
				"name": "DEV",
				"type": "go",
				"request": "launch",
				"mode": "debug",
				"program": "<absolute path to project>/main.go",
				"envFile": "<absolute path to project>/dev.env"
			},
			{
				"name": "SANDBOX",
				"type": "go",
				"request": "launch",
				"mode": "debug",
				"program": "<absolute path to project>/main.go",
				"envFile": "<absolute path to project>/sandbox.env"
			}
		]
	}
	```
1. Change the appropriate fields to match your dev environment
1. Add breakpoints, select your launch config from **Run and Debug**, and Start Debugging

> Note: Check out [this doc](https://github.com/golang/vscode-go/blob/master/docs/debugging.md) for more info on debugging Golang with VS Code

#### Enabling Webhook Service
1. Copy the k8s certificate to your `cert` folder in the operator using `kubectl cp zlifecycle-il-operator-system/<operator_pod>:/tmp/k8s-webhook-server/serving-certs <operator_project_root>/cert`
2. Make sure `ca.crt`, `tls.crt` and `tls.key` are in the `cert` folder
3. Add an environment variable in your `<environment>.env` file: `KUBERNETES_CERT_DIR=<operator_project_root>/cert`

## Other tools
1. Make sure all the variables in the env file start with `export <key>=<value`
2. Run `source <environment>.env`
3. Build the operator and run the executable
