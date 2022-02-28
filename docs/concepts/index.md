# Core Concepts

## Environment 

A logical grouping of all the Infrastructure Components that are needed to run business applications. The grouping includes components like networking, platform-eks, database, s3 buckets, and any other components.

## Components 

Logical groupings of 1 or more Infrastructure Resources that get provisioned together. For example, Networking is an Infrastructure Component with various Infrastructure resources like Virtual Private Cloud(VPC), Subnets, Internet Gateways, Route Tables, etc.

## Environment as Code

Environment as Code (EaC) is an abstraction over Infrastructure as Code that provides a declarative  way of defining an entire Environment. It has a Control Plane that manages the state of the environment, including relationships between various resources, Detects Drift as well enables Reconciliation. It also supports best practices  like Loose Coupling, Idempotency, Immutability, etc. for the entire environment. EaC allows teams to deliver entire environments rapidly and reliably, at scale.

To read more about this concept, go to [From Infrastructure as Code to Environment as Code](https://compuzest.com/2021/09/23/from-infrastructure-as-code-to-environment-as-code/) 