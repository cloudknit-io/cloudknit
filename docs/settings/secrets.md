# Environment Secrets

If you want to pass a secret to a component (for example a terraform module) you can create a secret using settings page and then reference it in the environment YAML. 

There are 3 scopes at which these secrets can be added, `org`, `team`, `environment` on the settings page.

![secret-scope](/assets/images/secret-scope.png "secret-scope")

To add a secret, go to the appropriate scope and click on the New button on the `Secrets` tab.

![environment-secrets](/assets/images/environment-secrets.png "environment-secrets")

Once added these secrets can be used in the secrets property of the environment YAML as below.

**Example**

```yaml
secrets:
  - name: bucket      # Terraform variable name
    key: s3-name      # Secret Id
    scope: org        # Should be one of: org, team, environment
```
