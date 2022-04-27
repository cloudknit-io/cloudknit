# Errors

Almost all errors are logged in the **concise logs** section of the [component details view](component-details-view.md). You can backtrack and get to the root cause of the error that happened during reconciliation.

Usually the errors that might occur are related to:-

- **AWS Secrets**: When secrets are not set and one is trying to provision an environment. See [secrets section](secrets.md) to set secrets. Another reason could be while deploying a particular resource for example:-

![AWS Error](../assets/images/aws-resource-error.png)
  
- **Terraform configuration**: Errors related to terraform usually occur because of **an error in the tfvars file** supplied by the user configured with an Environment yaml.

**Example:**

```
name            = "cust-team-env-vpc"
cidr            = "10.11.0.0/16"
azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
private_subnets = ["10.11.0.0/24", "10.11.1.0/24", "10.11.2.0/24"]
public_subnets  = ["10.11.200.0/24", "10.11.201.0/24", "10.11.202.0/24"]
enable_ipv6 = truee // Error: we added wrong value here

public_subnet_ipv6_prefixes  = [0, 1, 2]
private_subnet_ipv6_prefixes = [3, 4, 5]
public_subnet_tags = {
  "kubernetes.io/role/elb"      = 1
}
```

The above **tfvars** file will error out as `truee` is not a valid value.

- **Incorrect YAML**:
   - Formatting problem
   - Wrong properties
For all yaml issues a **notification** is shown on the UI, telling you the **problem part** in that yaml.

**Example**

```yaml
apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
  name: dev-checkout-sandbox
  namespace: zlifecycle
spec:
  teamName: checkout
  envName: sandbox  
  autoApprove: true
  component2: # Error Part
    - name: networking
      type: terraform
      module:
        source: aws
        name: vpc
      variablesFile:
        source: "git@github.com:zl-dev-tech/checkout-team-config.git"
        path: "sandbox/tfvars/networking.tfvars"
      outputs:
        - name: vpc_id
        - name: public_subnets
        - name: private_subnets
        - name: vpc_cidr_block
```

![yaml-notification](../assets/images/yaml-error.png "yaml notification")

In the above screenshot we can see that an unknown field **component2** is used in the **environment.spec** section of the yaml, which we can see in the **example**.