# AWS Credentials

![secrets](assets/images/secrets.png "Secrets")

You can **create** & **update** **AWS secrets** using zLifecycle secrets manager, which is accessible by **clicking** on the **Settings Navigation button**, as highlighted in the above image.

AWS section has three type of secrets:

![aws-secrets](assets/images/aws-secrets.png "aws-secrets")

* `Access Key Id`
* `Secret Access Key`
* `Session Token` [Optional]

These secrets are used by zlifecycle to provision your **AWS** environment. There are 3 scopes to which these secrets can be added, `company`, `team`, `environment`.

**NOTE:** You need to **add AWS Secrets before** provisioning an environment **for the first time**.

By default, zlifecycle tries to find secrets at environment level, then at team level and lastly at company level.
