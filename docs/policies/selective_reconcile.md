# Selective Reconcile

Selective reconcile can be used when only certain components should be re-provisioned within an environment. With this configuration re-provisioning an environment can leave some components untouched.

1. The First step is to identify the components you wish to selectively reconcile. Add **selectiveReconcile** to the **spec** scope in yaml.
```yaml
selectiveReconcile:
    tagName: helloWorldComponentType
    tagValues: [app, data]
```

**tagName**: This is a string property that should match the tags property of a component.

**tagValues**: An array of string values. Here we specify the values of the property type specified in `tagName`.

2. The Second part is the **tags** property, in the **component** scope.

```yaml
tags:
    - name: helloWorldComponentType 
      value: app
```



**name**: The name property needs to be exactly what we supplied in the **tagName** property of **selectiveReconcile** which in our example is `componentType`.
**value**: Value needs to be one of the strings passed in the **tagValues** property of **selectiveReconcile** which in our example is `app`.

**NOTE**: If tearing down or reconciling an environment, components left out of the selectiveReconcile will show a status of `Skipped Teardown` or `Skipped Reconcile`.

## Skip Mode

It is also possible to use selectiveReconcile to apply the inverse policy, and _skip_ components matching the tags, reconciling all others. 

```yaml
selectiveReconcile:
    skipMode: true
    tagName: helloWorldComponentType
    tagValues: [app, data]
```

In that case zLifecycle compares the above properties and sets up matching components to be skipped.

