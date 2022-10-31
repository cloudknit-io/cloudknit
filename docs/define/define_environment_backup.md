

<h4 id="auto-approve" style="font-weight: 200; letter-spacing: 2px">
  Auto Approve
</h4>

When you provision or teardown an **environment**, the terraform plan needs to be approved. The approval step can be automated, by setting the flag to `true`. By default, this flag is set to `false`, requiring the user to manually approve.


This property allows **{{ company_name }}** to skip the approval process.

**OPTIONAL**: defaulted to false if not provided

```yaml
autoApprove: true
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="teardown" style="font-weight: 200; letter-spacing: 2px">
  Teardown
</h4>

This property tells **{{ company_name }}** to destroy an environment, so if you are **provisioning** an environment **remember to either remove it or set it to false**

**OPTIONAL**: default value is false

You can find more information about teardown [here](/policies/teardown).

```yaml
teardown: true
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="selective-reconcile" style="font-weight: 200; letter-spacing: 2px">
  Selective Reconcile (Optional)
</h4>

This property tells **{{ company_name }}** to **skip** certain components based on **tagName** and **tagValues** properties.

**OPTIONAL**

You can find more information about **Selective Reconcile** [here](/policies/selective-reconcile).

**NOTE:** Works in conjunction with [**tags**](#component-tags) property of component.

```yaml
selectiveReconcile:
  tagName: string
  tagValues: [string, string]
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="spec-components" style="font-weight: 200; letter-spacing: 2px">
  Components
</h4>

This property contains an array of components that an environment is comprised of.

```yaml
components: []
```

See [components section](#components)

<div style="background-color: #ccc; height: 1px"></div>

<h4 style="font-weight: 200; letter-spacing: 2px">
  Usage
</h4>

```yaml
teamName: client
envName: demo
autoApprove: true
teardown: false
# Add if you want to skip components
selectiveReconcile:
  tagName: string
  tagValues: [string, string]
components: []
```

---

### Components

YAML Properties:-

  - [Name](#component-name)
  - [Type](#component-type)
  - [Destroy](#component-destroy)
  - [AWS Provider](#component-aws-provider)
  - [Modules](#component-modules)
  - [Outputs](#component-outputs)
  - [Variables](#component-variables)
  - [Secrets](#component-secrets)
  - [Overlay Data](#component-overlay-data)
  - [Overlay Files](#component-overlay-files)
  - [Depends On](#component-depends-on)
  - autoApprove (at component level)

This is the most intimidating part of your environment yaml file. Let's decipher it step by step.

<h4 id="component-name" style="font-weight: 200; letter-spacing: 2px">
  Name
</h4>

Name of the environment component

```yaml
name: static-assets
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-type" style="font-weight: 200; letter-spacing: 2px">
  Type
</h4>

Terraform is currently the only supported type

```yaml
type: terraform
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-destroy" style="font-weight: 200; letter-spacing: 2px">
  Destroy
</h4>

This property is used to override the `teardown` property of [spec](#spec), which applies to all components in the Environment file. It will override the `teardown` value for the component it is applied to.

**NOTE**: This property overrides the teardown property at the spec level, which means that if teardown is false and destroy is true, the current component gets destroyed.

**OPTIONAL**: Default value is false.

You can find more information about teardown [here](/policies/teardown)

```yaml
destroy: false
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-aws-provider" style="font-weight: 200; letter-spacing: 2px">
  AWS Provider (Optional)
</h4>

Below is an example portraying how to add an aws provider configuration to a component.

```yaml
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
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-modules" style="font-weight: 200; letter-spacing: 2px">
  Modules
</h4>

- [Public Module](#public-module)
- [Private Module](#private-module)
  <br/>

- <h5 id="public-module">Public Module</h5> 
    Currently only AWS modules are supported, which one can reference from https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest

  - <h6>Path (Optional)</h6>
    If the actual module is in a subdirectory (MonoRepo with multiple terraform modules), use `path` to specify the module

    ```yaml
    path: path/to/module
    ```

  <h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

  ```yaml
  module:
    source: aws
    name: s3-bucket
    path: path/to/module
  ```

- <h5 id="private-module">Private Module</h5>
    For private modules you need to specify full path to the terraform module.
    <h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>
    
    ```yaml
      module:
        source: "git@github.com:SebastianUA/terraform-aws-sagemaker"
    ```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-tags" style="font-weight: 200; letter-spacing: 2px">Tags (Optional)</h4>

Adds additional information to the component.

**Required**: When we are using [**selectiveReconcile**](#selective-reconcile) to skip components

**name**: For selective reconcile to work this needs to be the same value used in **tagName** property of **selectiveReconcile**

**value**: Value of the tag.

<h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

```yaml
tags:
  - name: helloWorldComponentType # for selective reconcile to work this needs to be the same value used in tagName property of selectiveReconcile
    value: data
  - name: cloudProvider
    value: aws
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-outputs" style="font-weight: 200; letter-spacing: 2px">Outputs</h4>

If the module supports outputs, name them here so they can be later referenced in `variables` block using `valueFrom`

<h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

```yaml
outputs:
  - name: bucket_arn
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-variables" style="font-weight: 200; letter-spacing: 2px">Variables (Optional)</h4>

```yaml
variables: []
```

Inline variables (will get injected into the terraform module when TF code is generated). This array is a combination of `name` and `value` or `valueFrom`.

  <br/>

- **Value**: String type.

- **ValueFrom**: reference an output defined in a previous module using [`outputs`](#module-outputs) block

<h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

```yaml
variables:
  - name: bucket
    value: "org-tech-client-demo-static-assets"
  - name: acl
    valueFrom: s3-common.bucket_acl
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-secrets"  style="font-weight: 200; letter-spacing: 2px">Secrets (Optional)</h4>

```yaml
secrets: []
```

References secret values which are added through the **{{ company_name }}** UI.

See [Secrets](/policies/secrets) Section.

  <br/>

- **Name**: Name of the terraform module variable.
- **Key**: Secret name entered in **{{ company_name }}** UI settings page.
- **Scope**: Scope configures secret scope granularity.

  - Org
  - Team
  - Environment

<h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

```yaml
secrets:
  - name: bucket
    key: s3-name
    scope: org
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-overlay-data" style="font-weight: 200; letter-spacing: 2px">Overlay Data (Optional)</h4>

```yaml
overlayData: []
```

Array of files to be generated and bundled with the environment component.

- **Name**: Name of the file.
- **Data**: Content of the file (generally it is a multi-line string).

<h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

```yaml
overlayData:
  - name: cloud-init.sh
    data: |
      #!/bin/sh
      echo "Starting cloud init"
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-overlay-files" style="font-weight: 200; letter-spacing: 2px">Overlay Files (Optional)</h4>

```yaml
overlayFiles: []
```

Array of external files which will be bundled with the environment component

- **Source**: Repo where the file is located.
- **Paths (array)**: Paths to files in the `source` repo. Types of paths:-
  - _Directory_
  - _Individual Overlay Files_
- **Ref (Optional)**: Reference to the branch, commit, head etc from where we will pull the file.

<h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

```yaml
overlayFiles:
  - source: "git@github.com:org-tech/client-config.git"
    paths:
      - demo/files/
      - demo/overlay.txt
      - demo/overlay2.txt
    ref: HEAD
```

<div style="background-color: #ccc; height: 1px"></div>

<h4 id="component-depends-on" style="font-weight: 200; letter-spacing: 2px">Depends On</h4>

```yaml
dependsOn: []
```

Add module dependencies (array of environment component names), which are to be resolved before the current component is processed.

Array includes [name](#component-name) property of the component.

<h5 style="font-weight: 200; letter-spacing: 2px">Usage</h5>

```yaml
dependsOn: [networking]
```

---

### Examples

- ##### networking

  ```yaml
  - name: networking
    type: terraform
    module:
      source: aws
      name: vpc
    variablesFile:
      source: "git@github.com:org-tech/client-config.git"
      path: "demo/tfvars/networking.tfvars"
    # if the module supports outputs, name them here so they can be later referenced in `variables` block using `valueFrom`
    outputs:
      - name: private_subnets
  ```

- ##### platform-eks

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
      source: "git@github.com:org-tech/client-config.git"
      # path to the file in the `source` repo
      path: "demo/tfvars/platform-eks.tfvars"
  ```

- ##### eks-addons

  ```yaml
  - name: eks-addons
    type: terraform
    dependsOn: [platform-eks]
    module:
      source: aws
      name: s3-bucket
    variablesFile:
      source: "git@github.com:org-tech/client-config.git"
      path: "demo/tfvars/eks-addons.tfvars"
  ```

- ##### platform-ec2

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
      source: "git@github.com:org-tech/client-config.git"
      path: "demo/tfvars/ec2.tfvars"
  ```

- ##### Full Fledged Yaml

  ```yaml
  apiVersion: stable.compuzest.com/v1
  kind: Environment
  metadata:
    name: org-tech-client-demo
    namespace: {{ company_name }}
  spec:
    teamName: client
    envName: demo
    # Use it to skip some components
    selectiveReconcile:
      tagName: string
      tagValues: [string, string]
  components:
    - name: static-assets
      type: terraform
      tags:
        - name: helloWorldComponentType # for selective reconcile to work this needs to be the same value used in tagName property of selectiveReconcile
          value: data
        - name: cloudProvider
          value: aws
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
          value: "org-tech-client-demo-static-assets"
        - name: acl
          valueFrom: s3-bucket.bucket_arn
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
        - source: "git@github.com:org-tech/client-config.git"
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
        source: "git@github.com:org-tech/client-config.git"
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
        source: "git@github.com:org-tech/client-config.git"
        path: "demo/tfvars/platform-eks.tfvars"
    - name: eks-addons
      type: terraform
      dependsOn: [platform-eks]
      module:
        source: aws
        name: s3-bucket
      variablesFile:
        source: "git@github.com:org-tech/client-config.git"
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
        source: "git@github.com:org-tech/client-config.git"
        path: "demo/tfvars/ec2.tfvars"
  ```

  ***
