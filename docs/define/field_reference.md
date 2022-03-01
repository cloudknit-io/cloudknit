# Field Reference

An environment YAML allows you to provide all the details of an environment.

## Sections

An environment YAML allows you to provide all the details of an environment. It has following main sections:

### Fields

| Field Name | Field Type | Description   |
|:----------:|:----------:|---------------|
|`apiVersion`|`string`|APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info [here](https://git.io.k8s.community/contributors/devel/sig-architecture/api-conventions.md#resources) |
|`kind`|`string`|Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info [here](https://git.io.k8s.community/contributors/devel/sig-architecture/api-conventions.md#types-kinds) |
|`metadata`|[`Metadata`](#Metadata)|_No description available_|
|`spec`|[`spec`](#spec)|_No description available_|

## Customer Resource Header

An environment YAML allows you to provide all the details of an environment. It has following main sections:

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
|`autoApprove`|`boolean`| To skip the manual approval step of applying the changes to a workflow, set this flag to `true`. If not set, it will default to `false` |
|`teardown`|`boolean`| To teardown an environment, set this flag to `true`. If you are creating a new environment, it must be omitted or set to `false`. If omitted, it will default to `false`. Environment teardown is composed as a 2-step process: First step is to update the `teardown` flag to `true` and wait for the environment to get destroyed (monitor progress on zLifecycle UI). Second step is to delete the Environment yaml from github for permanent deletion of an environment |
|`selectiveReconcile`| `array` | _No description available_ |
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
