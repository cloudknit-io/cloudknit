# Environment Secrets

If you want to pass a secret to a component (for example a terraform module) you can create a secret using settings page and then reference it in the environment YAML. 

Similar to AWS secrets we have 3 scopes here as well `company`, `team`, `environment`. To add a secret, go to the appropriate scope and click on the Add button.

![environment-secrets](../assets/images/environment-secrets.png "environment-secrets")

Once added these secrets can be used in the secrets property of the environment YAML as below.

**Example**
```yaml
secrets:
  - name: bucket
    key: s3-name # secret id
    scope: org
```
