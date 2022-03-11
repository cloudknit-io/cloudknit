# Example YAMLs

See below some examples of Environment YAML.

<details>
  <summary>Environment YAML with tfvars file</summary>
```yaml
apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
  name: dev-zmart-checkout-staging
  namespace: zmart-config
spec:
  teamName: checkout
  envName: staging
  autoApprove: true
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
        - name: vpc_id
        - name: public_subnets
        - name: private_subnets
      tags:
        - name: componentType
          value: app
    - name: platform-eks
      type: terraform
      dependsOn: [networking]
      module:
        source: git@github.com:terraform-aws-modules/terraform-aws-eks.git?ref=v17.24.0
      variablesFile:
        source: "https://github.com/zl-zmart-tech/checkout-team-config.git"
        path: "staging/tfvars/platform-eks.tfvars"
      overlayFiles:
        - source: "https://github.com/zl-zmart-tech/checkout-team-config.git"
          paths:
            - staging/tf/eks.tf
    - name: eks-addons
      type: terraform
      dependsOn: [platform-eks]
      module:
        source: aws
        name: s3-bucket
      outputs:
        - name: s3_bucket_arn
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
