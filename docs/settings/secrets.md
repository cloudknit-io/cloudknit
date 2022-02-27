# Secrets

![secrets](assets/images/secrets.png "Secrets")

An end-user can **create** and **update** their **secrets** using our secrets manager, which is accessible by **clicking** on the **Settings Navigation button**, as highlighted in the above image.

zLifecycle has two types of secrets:
* **AWS:** As the name suggests, helps you to create secrets for your AWS account.
   
* **Custom:** Custom secrets required during the reconciling of the environment.

### Custom Secrets

These are the secrets that a user can use at his own behest.

Similar to AWS secrets we have 3 scopes here as well `company`, `team`, `environment`.

To add a secret, select the scope and click on the Add button.

![custom-secrets](assets/images/custom-secrets.png "custom-secrets")

Here you can **provide** the **secret name, value** and you are set.

These secrets are used in the secrets property of the environment yaml.

**Example**
```yaml
secrets:
  - name: bucket
    key: s3-name # secret id
    scope: org
```
