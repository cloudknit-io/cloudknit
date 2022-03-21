# Selective Reconcile

Selective reconcile can be used when you don't need to re-provision the entire environment but only reconcile certain components.

1. The First step is to identify the components you wish to selectively reconcile. Add **selectiveReconcile** to the **spec** scope in yaml.
```yaml
selectiveReconcile:
    tagName: componentType
    tagValues: [app, data]
```

**tagName**: This is a string property that is used by tags property of component.

**tagValues**: An array of string values. Here we specify the values of the property type specified in `tagName`.

2. The Second part is the **tags** property, in the **component** scope.

```yaml
tags:
    - name: componentType
      value: app
```



**name**: The name property needs to be exactly what we supplied in the **tagName** property of **selectiveReconcile** which in our example is `componentType`.
**value**: Value needs to be one of the strings passed in the **tagValues** property of **selectiveReconcile** which in our example is `app`.

Now zLifecycle compares the above properties and sets up a component to be skipped.

**NOTE**: If we supply **teardown** as **true**, then the status would be `Skipped Teardown` else it is `Skipped Reconcile`.


