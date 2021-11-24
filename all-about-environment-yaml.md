# Environment YAML

An Environment YAML follows the following template:

- [Environment YAML](#environment-yaml)
    - [Metdata](#metdata)
    - [Spec](#spec)
    - [Components](#components)
      - [static-assets](#static-assets)
      - [networking](#networking)
      - [platform-eks](#platform-eks)
      - [eks-addons](#eks-addons)
      - [platform-ec2](#platform-ec2)
    - [Final YAML](#final-yaml)

Our YAML file always starts with following yaml:

```yaml
apiVersion: stable.compuzest.com/v1
kind: Environment
```

### Metdata

Metadata scope contains name of your environment which follows the following pattern `{orgName}-{teamName}-{environmentName}`

**orgName** is your organisation's name.
**teamName** and **environmentName** are provided in the [spec](#spec) scope.

```yaml
metadata:
  # Environment CRD k8s object name
  name: orgtech-cloudsync-demo
  # namespace is `zlifecycle` for every yaml you create
  namespace: zlifecycle
```

### Spec

Spec consists of following parameters:
`teamName`: `string`
`envName`: `string`
`autoApprove`: `boolean` [`OPTIONAL`]
`teardown`: `boolean` [`OPTIONAL`]
[`components`](#components): `Array`

```yaml
# name of the team CRD which this environment belongs to (also used to create metadata.name)
teamName: cloudsync

# name of the environment (used to create metadata.name)
envName: demo

# OPTIONAL: defaulted to false
# Skip approval process on zLifecycle UI
autoApprove: true

# OPTIONAL: use this flag for destroying an environment

# when creating a new environment, it must be omitted or set to `false`

# environment teardown is composed as a 2-step process
# Step 1. Update the `teardown` flag to `true` and wait for the environment to get destroyed (monitor progress on zLifecycle UI)
# Step 2. Delete the Environment yaml from github for permanent deletion of an environment

# NOTE: without the second step, you could use `teardown` to provision and restore temporary environments by toggling its value
teardown: false
components: [] # Array of components
```

### Components

This is the most intimidating part of your environment yaml file. Let's decipher it step by step.

A component consists of the following properties:

Properties such as `name`, `terraform` are self explanatory.

**Destroy** property takes in a boolean value, telling zlifecycle whether to destroy the environment or not. This property overrides **teardown** value provided at [spec](#spec) level, if it is set to `true`.

**DestroyProtection** property comes into use when we set the **teardown** [see [spec](#spec)] property to true. With this property set to true, it protects this component from being destroyed while tearing down an environment.

#### static-assets

```yaml
# array of environment components
components:
  # name of the environment component
  - name: static-assets

    # `terraform` is currently the only supported type
    type: terraform

    # OPTIONAL: same logic as global `teardown`, here it applies on environment component level (default is `false`)
    destroy: false

    # OPTIONAL: If set to true, do not run destroy actions against this component (default is `false`)
    destroyProtection: true

    # OPTIONAL: Configuration block for AWS provider
    aws:
      # OPTIONAL: AWS region
      region: us-east-1
      # OPTIONAL: Configuration block for AWS Assume Role
      assumeRole:
        # Role ARN which to assume
        roleArn: arn:aws:iam::account-id:role/zl-allow-assume-networking
        # OPTIONAL: External ID
        externalId: test_id1
        # OPTIONAL: Session Name
        sessionName: some_session


    # you can either reference a public or private module

    # for public AWS modules
    module:
      # `aws` is currently the only supported type
      source: aws

      # public terraform modules can be referenced here
      # this references https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest
      name: s3-bucket

      # OPTIONAL: if the actual module is in a subdirectory (monorepo with multiple terraform modules), use `path` to specify where is the module
      path: path/to/module

    # for private AWS modules
    module:
    # full path to the terraform module
    source: "git@github.com:SebastianUA/terraform-aws-sagemaker"

    # if the module supports outputs, name them here so they can be later referenced in `variables` block using `valueFrom`
    outputs:
      - name: bucket_arn

    # OPTIONAL: inline variables (will get injected into the terraform module when TF code is generated)
    variables:
      # array of `name -> value` objects
      - name: bucket
        value: "org-tech-cloudsync-demo-static-assets"
        # example of how to fetch a variable from `outputs`
      - name: acl
        # reference an output defined in a previous module using `outputs` block
        valueFrom: s3-common.bucket_acl


    # OPTIONAL: reference secret values which are added through the zLifecycle UI
    secrets:
      # name of the terraform module variable
      - name: bucket
        # secret name entered in zLifecycle UI settings page
        key: s3-name
        # scope configures secret scope granularity
        # valid scopes are org, team, environment and component
        scope: org


    # OPTIONAL: array of files to be generated and bundled with the environment component
    overlayData:
      # name of the file
      - name: cloud-init.sh
        # content of the file (generally it is a multi-line string)
        data: |
          #!/bin/sh
          echo "Starting cloud init"


    # OPTIONAL: external files which will be bundled with the environment component
    overlayFiles:
      # repo where the file is located
      - source: "git@github.com:org-tech/cloudsync-config.git"
        # paths to files in the `source` repo
        paths:
          # directory with overlay files
          - demo/files/
          # individual overlay files
          - demo/overlay.txt
          - demo/overlay2.txt
        # OPTIONAL: from which ref to pull the file (HEAD, branch name, commit SHA...)
        ref: HEAD
```

#### networking

```yaml
- name: networking
  type: terraform
  module:
    source: aws
    name: vpc
  variablesFile:
    source: "git@github.com:org-tech/cloudsync-config.git"
    path: "demo/tfvars/networking.tfvars"
  # if the module supports outputs, name them here so they can be later referenced in `variables` block using `valueFrom`
  outputs:
    - name: private_subnets
```

#### platform-eks

```yaml
- name: platform-eks
  type: terraform
  # add module dependencies (array of environment component names)
  dependsOn: [networking]
  module:
    source: aws
    name: s3-bucket
  # instead of inline variables, pass a tfvars file
  variablesFile:
    # repo where the file belongs
    source: "git@github.com:org-tech/cloudsync-config.git"
    # path to the file in the `source` repo
    path: "demo/tfvars/platform-eks.tfvars"
```

#### eks-addons

```yaml
- name: eks-addons
  type: terraform
  dependsOn: [platform-eks]
  module:
    source: aws
    name: s3-bucket
  variablesFile:
    source: "git@github.com:org-tech/cloudsync-config.git"
    path: "demo/tfvars/eks-addons.tfvars"
```

#### platform-ec2

```yaml
- name: platform-ec2
  type: terraform
  dependsOn: [networking]
  module:
    source: aws
    name: ec2-instance
  # example of using both inline variables and tfvars file
  variables:
    - name: subnet_id
      # example of how to fetch a variable from `outputs`
      valueFrom: networking.private_subnets[0]
  variablesFile:
    source: "git@github.com:org-tech/cloudsync-config.git"
    path: "demo/tfvars/ec2.tfvars"
```

### Final YAML

The final yaml would result into something similar like the example given below.

```yaml
apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
  name: org-tech-cloudsync-demo
  namespace: zlifecycle
spec:
  teamName: cloudsync
  envName: demo
  autoApprove: true # OPTIONAL
  teardown: false # OPTIONAL (Not required when provisioning environment)
  components:
    - name: static-assets
      type: terraform
      aws:
        region: us-east-1
        assumeRole:
          roleArn: arn:aws:iam::account-id:role/zl-allow-assume-networking
          externalId: test_id1
          sessionName: some_session
      module:
        source: aws
        name: s3-bucket
        path: path/to/module
      outputs:
        - name: bucket_arn

      variables:
        - name: bucket
          value: "org-tech-cloudsync-demo-static-assets"
        - name: acl
          valueFrom: s3-common.bucket_acl
      secrets:
        - name: bucket
          key: s3-name
          scope: org
      overlayData:
        - name: cloud-init.sh
          data: |
            #!/bin/sh
            echo "Starting cloud init"
      overlayFiles:
        - source: "git@github.com:org-tech/cloudsync-config.git"          
          paths:
            - demo/files/
            - demo/overlay.txt
            - demo/overlay2.txt
          ref: HEAD
    - name: networking
      type: terraform
      module:
        source: aws
        name: vpc
      variablesFile:
        source: "git@github.com:org-tech/cloudsync-config.git"
        path: "demo/tfvars/networking.tfvars"
      outputs:
        - name: private_subnets
    - name: platform-eks
      type: terraform
      dependsOn: [networking]
      module:
        source: aws
        name: s3-bucket
      variablesFile:
        source: "git@github.com:org-tech/cloudsync-config.git"
        path: "demo/tfvars/platform-eks.tfvars"
    - name: eks-addons
      type: terraform
      dependsOn: [platform-eks]
      module:
        source: aws
        name: s3-bucket
      variablesFile:
        source: "git@github.com:org-tech/cloudsync-config.git"
        path: "demo/tfvars/eks-addons.tfvars"
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
        source: "git@github.com:org-tech/cloudsync-config.git"
        path: "demo/tfvars/ec2.tfvars"
```
