# CloudKnit: An Open Source Solution for Managing Cloud Environments

**{{ company_name }}** is an open-source progressive delivery platform for managing cloud environments.

It enables organizations to Define entire environments in a declarative way, Provision them, Detect and Reconcile Drift, and Teardown environments when no longer needed. It also comes with dashboards to help visualize environments and observe them.

CloudKnit is based on a concept called [Environment as Code](https://www.cloudknit.io/blog/from-infrastructure-as-code-to-environment-as-code). Some people have started calling it Declarative Pipelines.

> *Note: We are not a big fan of using Pipeline and Declarative together as Pipeline to us means a sequence of steps which conflicts with what Declarative means.*

Environment as Code (EaC) is an abstraction over Cloud-native tools that provides a declarative way of defining an entire Environment. It has a Control Plane that manages the state of the environment, including resource dependencies, and drift detection and reconciliation.

![Where CloudKnit connects with existing tools](/assets/images/existing-tools.png)
*<center>Diagram 1: Where does CloudKnit fit in with existing tools</center>*
