# Example YAMLs

See below some examples of Environment YAML.

<details>
  <summary>Environment YAML with tfvars file</summary>
```yaml
apiVersion: stable.cloudknit.io/v1
kind: Environment
metadata:
  name: dev-zmart-checkout-staging
  namespace: zmart-config
spec:
  teamName: checkout
  envName: staging
  components:
    - name: networking
      type: terraform
      module:
        source: aws
        name: vpc
      variablesFile:
        source: "https://github.com/zl-zmart-tech/checkout-team-config.git"
        path: "staging/tfvars/networking.tfvars"
      outputs:
        - name: private_subnets
    - name: platform-ec2
      type: terraform
      dependsOn: [networking]
      module:
        source: aws
        name: ec2-instance
      variables:
        - name: subnet_id
          valueFrom: networking.private_subnets[0]
      variablesFile:
        source: "https://github.com/zl-zmart-tech/checkout-team-config.git"
        path: "staging/tfvars/ec2.tfvars"
```
</details>

<details>
  <summary>Networking tfvars file</summary>
```YAML
name            = "dev-zmart-checkout-staging-vpc"
cidr            = "10.12.0.0/16"
azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
private_subnets = ["10.12.0.0/19", "10.12.64.0/19", "10.12.128.0/19"]
public_subnets  = ["10.12.48.0/20", "10.12.112.0/20", "10.12.176.0/20"]
enable_nat_gateway                          = true
single_nat_gateway                          = true
enable_dns_hostnames                        = true

tags = {
  Terraform                                 = "true"
  Environment                               = "dev-zmart"
  "kubernetes.io/cluster/dev-checkout-staging-eks"           = "shared"
}

public_subnet_tags = {
  Terraform                                 = "true"
  Environment                               = "dev-zmart"
  "kubernetes.io/cluster/dev-checkout-staging-eks"           = "shared"
  "kubernetes.io/role/elb"                  = 1
}

private_subnet_tags = {
  Terraform                                 = "true"
  Environment                               = "dev-zmart"
  "kubernetes.io/cluster/dev-checkout-staging-eks"           = "shared"
  "kubernetes.io/role/internal-elb"         = 1
}

database_subnet_tags = {
  Terraform                                 = "true"
  Environment                               = "dev-zmart"
}
```
</details>


<details>
  <summary>EC2 tfvars file</summary>
```YAML
name            = "dev-zlab-checkout-staging-platform-ec2"
ami             = "ami-06aac012759a08cec"
instance_type   = "t3.medium"
```
</details>
