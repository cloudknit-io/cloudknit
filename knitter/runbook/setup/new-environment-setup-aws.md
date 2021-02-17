# Setup a new Environment on AWS 

## Overview

Steps to setup a new Environment on AWS 

## When to use this runbook
This is to be used when you are setting up brand new environment on AWS 

## Initial Steps Overview

1. [Create tfvars files in terraform-provisioner](#create-tfvars-files-in-terraform-provisioner)
1. [Other Changes](#other-changes)

## Detailed Steps

#### Create tfvars files in terraform-provisioner
For all Terraform Configs create tfvars file with the environment name

- For example: under `zlifecycle-provisioner/aws-vpc/tfvars` folder create `dev.tfvars` by cloning `sandbox.tfvars` file and then making appropriate changes
- Do the same for `aws-eks` and `k8s-addons` terraform configs

#### Other Changes
- Update `bootstrap/bootstrap_zlifecycle.sh` file to add a new environment in the selection options
- In `bootstrap/aws` folder create a company config file by cloning an existing company config file and changing it appropriately

