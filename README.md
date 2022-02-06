# zlifecycle-user-docs

# Index

* Overview
    * What is zLifecycle
    * Why zLifecycle
    * zLifecycle Documentation
        * Getting Started
        * Define
        * Provision
        * Teardown
* Core Concepts
* Getting on-boarded on to zLifecycle
* User Guide
    * Getting Started
        * Setting up zLifecycle
        * [Onboarding a Team](team/onboarding.md)
        * [Provisioning your First Environment](provision/first-environment.md)
        * [Tearing down your First Environment](teardown/first-environment.md)
    * [Define Environment](define/define-environment.md)
    * Provision Environment
    * Reconcile Environment
    * Teardown Environment
    * Dashboards

# Overview

## What is zLifecycle?

![image](https://user-images.githubusercontent.com/47644789/147984939-738f7535-be82-41ab-8f35-e684f8cdb3c7.png)

zLifecycle is a declarative, GitOps continuous delivery tool for Environments.

zLifecycle makes Environment Management easy across various cloud providers. It enables organizations to define, provision, detect drift + reconcile, and teardown environments across multiple cloud providers as well as on-premises. It also provides dashboards and reports for data-driven decision making.


## Why zLifecycle?

Environment definitions and configurations should be declarative and version controlled. Environment deployment and lifecycle management should be automated, auditable, and easy to understand.

## zLifecycle Documentation

### Getting Started

### Define

### Provision

### Teardown

# Core Concepts

## Teams

## Environment 

A logical grouping of all the Infrastructure Components that are needed to run business applications. The grouping includes components like networking, platform-eks, database, s3 buckets, and any other components.

## Components 

Logical groupings of 1 or more Infrastructure Resources that get provisioned together. For example, Networking is an Infrastructure Component with various Infrastructure resources like Virtual Private Cloud(VPC), Subnets, Internet Gateways, Route Tables, etc.



## Custom Resource Definition (CRD)
Environment YAML 

## Provisioning 

## Drift Detection + Reconciling

## Teardown

## 

# User Guide

## Reference YAML

## Field Reference

## Cost Calculation

## Selective Reconcile

This property tells zlifecycle to **skip** certain components based on **tagName** and **tagValues** properties. This is an optional field.

**NOTE:** Works in conjunction with [**tags**](#component-tags) property of component.

```yaml
selectiveReconcile:
  tagName: string
  tagValues: [string, string]
```

## Manual Reconcile

## Status
List of statuses + description

## Approval

When you provision or teardown an **environment**, the terraform plan needs to be approved. This is an optional property called `autoApprove` which can be added at spec level or component level.
The approval step can be automated, by setting the flag to `true`. By default, this flag is set to `false`, requiring the user to manually approve.

```yaml
autoApprove: true
```

## Team Onboarding

## Team member Onboarding

## Secrets Management
Set up AWS credentials
Set up other secrets

## Overlay Files

## AWS Provider

## Destroy Protection

## Possible Errors



* [How to create an Environment YAML file](all-about-environment-yaml.md)
* How to onboard a team
* [How to approve](approval.md)
* [How to add a Secret](secrets.md)
* [How to use Reconcile](reconcile.md)
* [Possible Errors](errors.md)
