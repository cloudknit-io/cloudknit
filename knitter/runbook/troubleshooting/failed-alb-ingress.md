# Environment Components (Terraform Configs)

## Overview

Steps to fix a known issue that occurs when recreating an EKS cluster and the ALB controller is in a failing state.

## When to use this runbook
When the load balancer url is not appearing on the alb `Ingress` resource in the `kube-system` namespace. `k get ingress -n kube-system` should show the load balancer domain.

## Initial Steps Overview

1. Confirm this is indeed the known issue by checking the logs of the alb pod of the `aws-load-balancer-controller` deployment in `kube-system`
1. Re-apply the `Ingress` resource in `kube-system`, either manually or via terraform where it is applied in its helm chart [here](https://github.com/CompuZest/helm-charts/blob/main/charts/zlifecycle-ingress/templates/ingress-alb.yaml)
