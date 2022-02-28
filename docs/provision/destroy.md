# Destroying Components

This property is similar to teardown property of [spec scope](define_environment/#spec), the only difference being, it applies on environment component level.

**NOTE**: This property overrides the teardown property at the spec level, which means that if teardown is false and destroy is true, the current component gets destroyed.

**OPTIONAL**: Default value is false.

You can find more information about teardown [here](link to the teardown page)
