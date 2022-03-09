# Onboard a New Team

When you want to onboard a new team use the following steps:

1. Create a repository that will consist of environments for the team
2. Make sure the zLifecycle github app has access to the new repo
3. Add team config yaml (like below) in the `company-config` repo

```yaml
apiVersion: stable.compuzest.com/v1
kind: Team
metadata:
  name: zmart-checkout-team
  namespace: zmart-config
spec:
  teamName: zmart-checkout-team
  configRepo:
    source: git@github.com:zmart-tech/zmart-checkout-team-config.git
    path: "."
```

4. Once you commit and push the change it will register the team repo with zLifecycle and watch for any updates
