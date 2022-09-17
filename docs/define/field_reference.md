# Field Reference

An environment YAML allows you to provide all the details of an environment.

## Sections

It has following main sections:

### Fields

Since Environment definition uses a Kubernetes Custom Resource the top section of the definition in YAML needs to follow its convention.

| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`apiVersion`|`string`| Custom Resource Header. Value needs to be `stable.compuzest.com/v1` |
|`kind`|`string`| Custom Resource Header. Value needs to be `Environment` |
|`metadata`|[`Metadata`](#Metadata)| Metadata about the Environment |
|`spec`|[`spec`](#spec)| Details about the Environment  |

## Custom Resource Header

```yaml
apiVersion: stable.compuzest.com/v1
kind: Environment
```

## Metadata

### Fields
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`name`|`string`| The name should be unique for every environment. To ensure that we follow below naming convention:- </br> `{company}-{team}-{environment}` </br></br>  **company** is your company's name. **team** and **environment** are defined in the [spec](#spec) section below |
|`namespace`|`string`| Namespace should be `{company}-config` |

---
<div style="background-color: #ccc; height: 1px"></div>

<h4 id="metadata-example" style="font-weight: 200; letter-spacing: 2px;">
  Example
</h4>

```yaml
metadata:
  name: zmart-checkout-dev
  namespace: zmart-config
```
---


## Spec

Spec contains details of the environment to be provisioned.

### Fields
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`teamName`|`string`| Name of the team to which this environment belongs |
|`envName`|`string`| Name of the environment |
|`autoApprove`|`boolean`| To skip the manual approval step of applying the changes to a workflow, set this flag to `true`. Default value is `false`. More info [here](https://docs.zlifecycle.com/define/approval/) |
|`teardown`|`boolean`| To teardown an environment, set this flag to `true`. Default value is `false`. More info [here](https://docs.zlifecycle.com/teardown/) |
|[`selectiveReconcile`](#selective-reconcile)| `array` | More info [here](https://docs.zlifecycle.com/define/selective_reconcile/) |
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
|`destroyProtection`|`boolean`| Optional field. If set to `true`, {{ company_name }} will not destroy this component (default is `false`) |
|`dependsOn`|`array`| Optional field. Array of environment component names, which this module depends on |
|[`secrets`](#secrets)|| This section references the secret values which are input through the {{ company_name }} UI |
|[`tags`](#tags)|| Tags are labels attached to components for the purpose of identification. It is an `array` of `string`  |
|[`variables`](#inline-variables)|| **Inline** variables, these will get injected into the terraform module when TF code is generated. `array` of `name -> value` objects |
|[`variablesFile`](#variables-from-a-file)|`string`| Variables can also be passed from an output defined in a previous module using `outputs` block, via a tfvars file |
|[`module`](#module)|`string`| Modules are containers for multiple resources that are used together. You can either reference a public or private module. |
|[`outputs`](#outputs)|`string`| _Output values make information about your infrastructure available on the command line, and can expose information for other components to use_. Output values are similar to return values in programming languages. If the module supports outputs, name them here so they can be referenced in `variables` block using `valueFrom` |
|[`overlayFiles`](#overlay-files)| | A file that contains additional information about the current items. By using an overlay file, the metadata of these items can be extended. |
|[`overlayData`](#overlay-data)| | Rather than have information pointing to the file with overlay information, you can also specify the data. |


<div style="background-color: #ccc; height: 1px"></div>
<div style="background-color: #ccc; height: 1px"></div>


### Selective Reconcile
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
| `skipMode` |`boolean`|Flag indicating if the tags mentioned for selective reconcile are to skip reconciliation. By default this is set to `false`. If you wish to skip reconciling some components, then set it to `true` and tag the components appropriately|
|`tagName`|`string`| Required, if using `selectiveReconcile`|
|`tagValues`|`array` of `string`| Required, if using `selectiveReconcile` |



### AWS
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`aws`| | Optional field. Configuration block for AWS provider. More info coming soon |


### Secrets
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`secrets`|| This section references the secret values which are input through the {{ company_name }} UI |
|`name`|`string`| Name of the terraform module variable |
|`key`|`string`| Secret name entered in {{ company_name }} UI settings page |
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
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`name`|`string`| Public terraform modules can be referenced [here](https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest). For private module, specify the full path |
|`path`|`string`| If the module is in a subdirectory (monorepo with multiple terraform modules), use this to specify the `path` |
|`source`|`string`| Required field. Currently `aws` is the only supported type|
|`version`|`string`| _No description available_ |

### Outputs
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`name`|`string`| Required field, if using `outputs` |
|`sensitive`|`boolean`| Optional field. Flag to indicate if the `output` is of sensitive nature. By default the value is set to `false`. To not display it in plaintext, set it to `true` |

### Overlay Files
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`source`|`array` of `string`| Required field, if using `overlayFiles`. Repo where the file is located |
|`path`|`array` of `string`| Required field, if using `overlayFiles`. Path to the file in the `source` repo |
|`ref`| | _No description available_ |

### Overlay Data
| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`data`|`string`| Required field, if using `overlayData`. Content of the file (generally it is a multi-line string) |
|`name`|`string`| Required field, if using `overlayData`. Name of the file, containing the afore specified data |


<div style="background-color: #ccc; height: 1px"></div>
