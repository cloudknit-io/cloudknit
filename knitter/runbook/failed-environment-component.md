
# Environment Components (Terraform Configs)

## Overview

Steps to change terraform state of customer environments

## When to use this runbook
When an environment cannot be synced or approved into a healthy state and requires troubleshooting or something like `terraform force unlock`

## Initial Steps Overview

1. Find failing workflow
1. Exec into the workflow pod 

## Other tactics
1. Manually unlock state in AWS console
1. Run workflow container as an image