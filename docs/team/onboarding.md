# Onboarding a New Team

When you want to onboard a new team use the following steps:
1. [Optional] Create a repository that will consist of environments for the team
2. [Optional] Give github svc account (this will be github apps soon) Admin access (this will be read & webhook creation perms with gh apps) 
2. Add team config yaml (like below) in the `zlifecycle-config` repo

```yaml
apiVersion: stable.compuzest.com/v1
kind: Team
metadata:
  name: zmart-checkout-team
  namespace: zlifecycle-config
spec:
  teamName: zmart-checkout-team
  configRepo:
    source: git@github.com:zmart-tech/zmart-checkout-team-config.git
    path: "."
```

3. Once you commit and push the change it will register the team repo with zLifecycle and watch for any updates
4. [Optional] Set the AWS credentials for all environments within the team on the Settings page
