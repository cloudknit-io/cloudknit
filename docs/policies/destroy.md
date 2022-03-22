# Destroying Components

This property is used to destroy a component. When an **environment** is brought down, we refer to it as [`teardown`](teardown.md) and when a **component** is brought down, it is `destroy`


**NOTE**: The component-level `destroy` property overrides the spec-level `teardown` property, which means that if `teardown` is `false` and `destroy` is `true`, the current component will be destroyed.

This is an optional field, with the default value as `false`.

You can find more information about teardown [here](teardown.md).
