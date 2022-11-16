# Core Concepts

## Environment

A logical grouping of all the Infrastructure Components that are needed to run business applications. The grouping includes components like networking, platform-eks, database, s3 buckets, and any other components.

## Components

Logical groupings of 1 or more Infrastructure Resources that get provisioned together. For example, Networking is an Infrastructure Component with various Infrastructure resources like Virtual Private Cloud(VPC), Subnets, Internet Gateways, Route Tables, etc.

## Environment as Code

Environment as Code (EaC) is an abstraction over Infrastructure as Code that provides a declarative  way of defining an entire Environment. It has a Control Plane that manages the state of the environment, including relationships between various resources, Detects Drift as well enables Reconciliation. It also supports best practices  like Loose Coupling, Idempotency, Immutability, etc. for the entire environment. EaC allows teams to deliver entire environments rapidly and reliably, at scale.

To read more about this concept, go to [From Infrastructure as Code to Environment as Code](https://www.cloudknit.io/blog/from-infrastructure-as-code-to-environment-as-code)

## GitOps

GitOps extends Infrastructure as Code (IaC) and adds a workflow (Pull Request Process) to apply a change to the Production or any environment for that matter. It could also have a control loop that verifies periodically that the actual state of the infrastructure is the same as the desired state.

```
GitOps = IaC + (Workflow + Control Loop)
```

To read more about GitOps, go to [Infrastructure as Code: Principles, Patterns, and Practices](https://www.cloudknit.io/blog/principles-patterns-and-practices-for-effective-infrastructure-as-code) article & checkout the GitOps section.
