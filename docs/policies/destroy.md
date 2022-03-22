# Destroying Components

This property is used to destroy a component. When an **environment** is brought down, we refer to it as [`teardown`](teardown.md) and when a **component** is brought down, it is `destroy`


**NOTE**: The component-level `destroy` property overrides the spec-level `teardown` property, which means that if `teardown` is `false` and `destroy` is `true`, the current component will be destroyed.

This is an optional field, with the default value as `false`.

You can find more information about teardown [here](teardown.md).

## Destroy Protection

You can always safeguard a component from getting destroyed by applying some protection to it using the `destroyProtection` flag. When set to `true` the current component will not be destroyed, overriding all other flags including the environment-level `teardown`.

This is also an optional field, with the default value as `false`.

```
 components:
    - name: static-assets
      type: terraform
      destroy: false
      destroyProtection: true
```