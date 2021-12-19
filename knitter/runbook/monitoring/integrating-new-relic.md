# Integrate New Relic

## Overview

Steps to connect k8s cluster and services (APM) with New Relic

## When to use this runbook
This is to be used when you are setting up monitoring & observability for a new zLifecycle instance

## Initial Steps Overview

1. [Setting up slack webhook notification for Argo Workflows](#setting-up-slack-webhook-notification-for-argo-workflows)
1. [Installing & Configuring New Relic](#installing--configuring-new-relic)

## Detailed Steps

### Setting up slack webhook notification for Argo Workflows

Environment variable `env.SLACK_WEBHOOK_URL` must be provided to IL operator so it can correctly send alerts for failed workflows.

The preferred way of setting the correct variables is by setting the `slack_webhook_url` in the `terraform-modules/zl-app-addons` module by either using tfvars file or fetching it through `data` resource from AWS SSM Parameter Store

### Installing & Configuring New Relic

#### Obtaining the New Relic API key

1. Open New Relic UI
1. Click on the Avatar (Profile) icon in top right -> API Keys
1. Find the key named `Original account license key` -> click the `...` on the right -> Copy Key

#### Installing New Relic Cluster Agent

Terraform variables `install_new_relic` should be set to `"true"` and `new_relic_api_key` should be set to the API (license) key.

#### Enabling APM for IL operator and state manager

Environment variables `env.ENABLE_NEW_RELIC` should be set to `"true"` and `env.NEW_RELIC_API_KEY` should be set to the API (license) key in the Helm charts

The preferred way of setting the correct variables is by setting the `enable_new_relic` and `new_relic_api_key` terraform variables in the `terraform-modules/zl-app-addons` module by either using tfvars file or fetching it through `data` resource from AWS SSM Parameter Store

#### Configuring alerts

##### Create notification channel

1. Open New Relic UI
1. Select `Alerts & AI` from the top nav bar -> Select `Channels` on the left sidebar -> New notification channel
1. Select `Slack` as channel type and connect it to a specific slack channel which will receive alerts (the UI is pretty intuitive)

##### Configuring alerts

1. Open New Relic UI
1. Select `More` from the top nav bar -> Workload Views -> Create a Workload
1. Name it using the format <env>-<service>
1. Select the appropriate `Service - APM` entity, Deployment entity and Kubernetes cluster entity (example for zlifecycle-il-operator: zlifecycle-il-operator Service - APM, zlifecycle-il-operator Kubernetes deployment and dev-eks Kubernetes cluster)
1. Click on `Create a workload`
1. Select the `Errors Inbox` from the top nav bar -> Select the newly created workload -> Click `Set Up Notification` -> Link it with your notification channel
