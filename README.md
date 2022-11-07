# CloudKnit: An Open Source Cloud Environment Manager

[CloudKnit](https://github.com/cloudknit-io/cloudknit) is an open-source "progressive delivery platform" for managing cloud environments. It enables organizations to ***Define** entire environments in a declarative way, **Provision** them, **Detect** and **Reconcile** Drift, and **Teardown** environments when no longer needed. It also comes with dashboards to help visualize environments and observe them.

CloudKnit is based on a concept called [Environment as Code](https://www.zlifecycle.com/blog/from-infrastructure-as-code-to-environment-as-code). Some people have started calling it Declarative Pipelines.

> *Note: We are not a big fan of using Pipeline and Declarative together as Pipeline to us means a sequence of steps which is conflicts with what Declarative means.*

Environment as Code (EaC) is an abstraction over Cloud-native tools that provides a declarative way of defining an entire Environment. It has a Control Plane that manages the state of the environment, including relationships between various resources, and Detects Drift, and does Reconciliation.

![Where CloudKnit connects with existing tools](/assets/images/existing-tools.png)
*<center>Diagram 1: Where does CloudKnit fit in with existing tools</center>*

## Why we built CloudKnit

There are tools today that allow us to manage cloud environments but as the environments become more complex and teams look for advanced use cases like visualizing environments, replicating environments, promoting changes across environments, blue/green deployments, etc., existing tools fall short. 

This causes some teams to build and maintain in-house solutions. We want to make it easy for devops teams to manage complex environments, get to advanced cases, and to gain a competitive advantage.

Existing cloud-native tools like Terraform, Pulumi, Helm, ArgoCD, etc. are great at automating simple environments, but if you want an environment like the one below (Diagram 2) with infrastructure resources and cloud-native applications, your options are:

### Option 1: Monolith Infrastructure as Code & Application Deployment

Monolith IaC & Application deployments work well when environments are simple, but as things get complex, it becomes a nightmare to maintain.

### Option 2: Use loosely coupled Infrastructure as Code & Application Deployment & hand-roll pipelines to run them

Creating loosely-coupled components like networking, eks, rds, k8s-apps, etc. makes the individual components easier to manage, but the pipelines that run those components become complex. Pipeline code needs to manage the logic to run the various components in the correct order and handle various scenarios like failures and tearing down an entire environment for ephemeral environments.

Pipeline code is imperative, and users have to write logic on “HOW” to get an entire environment. This causes a maintenance nightmare. We have seen teams write hundreds of lines of pipeline code to specify the “HOW” to run various components to get an entire environment and it becomes unmanageable at one point.

![Where does CloudKnit fit in with existing tools](/assets/images/environment.jpeg)
*<center>Diagram 2: Example Environment</center>*

## Other Challenges

### Environment Replication is a pain

### Not easy to Visualize/Understand Environments 

### Drift Detection for the entire environment is difficult

### Not straightforward to Promote changes across environments

## How does CloudKnit work?

![CloudKnit](/assets/images/cloudknit.jpeg)
*<center>Diagram 3: CloudKnit</center>*

Environment management with CloudKnit is divided into 4 stages:

### Define

This stage allows you to define an entire environment. See example below: 

<details>
  <summary>Environment Definition</summary>

```

```
</details>

### Provision

Control Plane in Kubernetes using the Definition generates argo workflow YAMLs that then is used to run Terraform/Helm etc. in the right order. It also provides Visibility & Workflow.

### Detect Drift + Reconciliation

Like Kubernetes does drift detection for k8s apps & reconciles them to match the desired state in source control, EaC does drift detection for the entire environment (infra + apps) & reconciles them ((Preferably with Approval step for infra that shows the plan). 

### Teardown

Environment Teardown is pretty straightforward and can be done easily by setting the `teardown` flag to true. The Control plane generates the Argo workflow that runs destroy on various components in the correct order.

## Conclusion

We hope that by open-sourcing CloseKnit early, we can form a close-knit open-source community around it to make managing complex cloud environments easy.

For a deeper dive into CloudKnit, see the [architecture document](TBD), our [documentation](https://docs.cloudknit.io), and the [GitHub repo](https://github.com/cloudknit-io/cloudknit).

#### Terminologies

*Components: A logical grouping of 1 or more Infrastructure Resources or Applications that get provisioned together. For example, Networking is an Infrastructure Component with various Infrastructure resources like Virtual Private Cloud(VPC), Subnets, Internet Gateways, Route Tables, etc.*

*Environment: A logical grouping of all the Components needed to run business applications. The grouping includes components like networking, eks, database, k8s apps, etc.*