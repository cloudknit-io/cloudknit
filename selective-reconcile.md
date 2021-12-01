# Selective Reconcile

Selective reconcile is used when a user wants to **skip reconciliation** of a component be it teardown or provisioning.

To skip a component you need to follow following steps:-

- The First step is to add **selectiveReconcile** to the **spec** scope in yaml.
```yaml
selectiveReconcile:
    tagName: componentType
    tagValues: [app, data]
```

**tagName**: This is a string property that is used by tags property of component.

**tagValues**: This is array of string. Here we add tagValues, this is also used to validate the component that needs to be skipped.

- The Second part is the **tags** property of component.

```yaml
tags:
    - name: componentType
      value: app
```

**name**: The name property needs to be exactly what we supplied in the **tagName** property of **selectiveReconcile** which in our example is `componentType`.
**value**: Value needs to be one of the strings passed in the **tagValues** property of **selectiveReconcile** which in our example is `app`.

Now zLifecycle compares the above properties and sets up a component to be skipped.

**NOTE**: If we supply **teardown** as **true**, then the status would be `Skipped Teardown` else it is `Skipped Reconcile`.


