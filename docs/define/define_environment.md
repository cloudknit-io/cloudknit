# Define Environment

As the name says the `Define` step of the Environment Lifecycle allows you to define and entire entire environment. Once you create the environment definition, you will have to commit and push the changes to team repository. **{{ company_name }}** will automatically pickup the changes and provision/update/teardown the environment.

Environment definition uses a Kubernetes [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) which is in YAML format and provides a declarative way of defining an environment.

Check the [Field Reference](./field_reference.md) page for information about all the fields.
