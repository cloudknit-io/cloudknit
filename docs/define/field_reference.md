# Field Reference

An environment YAML allows you to provide all the details of an environment.

## Sections

It has following main sections:

### Fields

| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`apiVersion`|`string`|APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info [here](https://git.io.k8s.community/contributors/devel/sig-architecture/api-conventions.md#resources) |
|`kind`|`string`|Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info [here](https://git.io.k8s.community/contributors/devel/sig-architecture/api-conventions.md#types-kinds) |
|`metadata`|[`Metadata`](#Metadata)|_No description available_|
|`spec`|[`spec`](#spec)|_No description available_|

## Custom Resource Header

```yaml
apiVersion: stable.compuzest.com/v1
kind: Environment
```

## Metadata

### Fields
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`name`|`string`| The name should be unique for every environment, to ensure that we follow below naming convention:- `{company}-{team}-{environment}` For example: `zmart-checkout-dev` **company** is your company's name. **team** and **environment** are defined in the [spec](#spec) section below |
|`namespace`|`string`| Namespace should be `{Company}-config`. So if you company name used in zLifecycle is `zmart` then the namespace should be `zmart-config` |

---
<div style="background-color: #ccc; height: 1px"></div>

<h4 id="metadata-example" style="font-weight: 200; letter-spacing: 2px;">
  Example
</h4>

```yaml
metadata:
  # Environment Custom Resource name in k8s
  name: zmart-checkout-dev
  # namespace should be `zlifecycle` for all environments
  namespace: zlifecycle
```
---


## Spec

Spec contains the information about details of the environment to be provisioned.

### Fields
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`teamName`|`string`| Name of the team to which this environment belongs (also used to create [metadata.name](#metdata)) |
|`envName`|`string`| Name of the environment |
|`autoApprove`|`boolean`| To skip the manual approval step of applying the changes to a workflow, set this flag to `true`. If not set, it will default to `false`. More info [here](https://docs.zlifecycle.com/define/approval/) |
|`teardown`|`boolean`| To teardown an environment, set this flag to `true`. If you are creating a new environment, it must be omitted or set to `false`. If omitted, it will default to `false`. Environment teardown is composed as a 2-step process: First step is to update the `teardown` flag to `true` and wait for the environment to get destroyed (monitor progress on zLifecycle UI). Second step is to delete the Environment yaml from github for permanent deletion of an environment. More info [here](INSERT TEARDOWN LINK HERE) |
|`selectiveReconcile`| `array` | More info [here](https://docs.zlifecycle.com/define/selective_reconcile/) |
|`components`|`array`| Array of environment components |

---
<div style="background-color: #ccc; height: 1px"></div>

<h4 id="metadata-example" style="font-weight: 200; letter-spacing: 2px;">
  Example
</h4>

```yaml
spec:
  teamName: checkout
  envName: demo
```
---

## Components

Array of environment components.

### Fields
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`name`|`string`| Name of the environment component |
|`type`|`string`| `terraform` is currently the only supported type |
|`destroy`|`boolean`| Optional field. Flag for destroying a component. Default is `false`. More info [here](../destroy.md) |
|`destroyProtection`|`boolean`| Optional field. If set to `true`, zLifecycle will not destroy this component (default is `false`) |
|`dependsOn`|`array`| Optional field. Array of environment component names, which this module depends on |
|[`secrets`](#secrets)|| This section references the secret values which are input through the zLifecycle UI |
|[`tags`](#tags)|| Tags are labels attached to components for the purpose of identification. It is an `array` of `string`  |
|[`variables`](#inline-variables)|| **Inline** variables, these will get injected into the terraform module when TF code is generated. `array` of `name -> value` objects |
|[`variablesFile`](#variables-from-a-file)|`string`| Variables can also be passed from an output defined in a previous module using `outputs` block, via a tfvars file |
|[`module`](#module)|`string`| Modules are containers for multiple resources that are used together. You can either reference a public or private module. |

<div style="background-color: #ccc; height: 1px"></div>
<div style="background-color: #ccc; height: 1px"></div>



### AWS
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`aws`| | Optional field. Configuration block for AWS provider. More info coming soon |


### Secrets
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`secrets`|| This section references the secret values which are input through the zLifecycle UI |
|`name`|`string`| Name of the terraform module variable |
|`key`|`string`| Secret name entered in zLifecycle UI settings page |
|`scope`|`string`| Refers to what scope the secret is valid in. Valid scopes are `org`, `team` and `environment` |

### Tags
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`tags`|| Tags are labels attached to components for the purpose of identification. It is an `array` of `string`  |
|`name`|`string`| Type of tag |
|`value`|`string`| Identifying tags |

### Variables

#### Inline Variables
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`variables`|| **Inline** variables, these will get injected into the terraform module when TF code is generated. `array` of `name -> value` objects |
|`name`|`string`| Name of the variable |
|`value`|`string`| Value of the variable |
|`valueFrom`|`string`|  |

#### Variables from a file
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`variablesFile`|`string`| Variables can also be passed from an output defined in a previous module using `outputs` block, via a tfvars file |
|`path`|`string`| Required field, if using `variablesFile`. Path to the file in the `source` repo. |
|`ref`|`string`| _No description available_ |
|`source`|`string`| Required field, if using `variablesFile`. Repo where the variables file can be found. |

### Module

#### Public Module
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`name`|`string`| Public terraform modules can be referenced [here](https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest) |
|`path`|`string`| If the module is in a subdirectory (monorepo with multiple terraform modules), use this to specify the `path` |
|`source`|`string`| Required field. Currently `aws` is the only supported type|
|`version`|`string`| _No description available_ |


<div style="background-color: #ccc; height: 1px"></div>