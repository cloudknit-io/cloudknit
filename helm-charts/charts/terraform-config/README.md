Terraform Config

This helm chart is 
The ArgoCD Applications that represent a terraform conig (e.g. Networking layer) for a customer use this helm chart to deploy the control loop responsible for keeping that terraform config in sync. 

The ArgoCD Applications themselves are configured [here](https://github.com/CompuZest/zlifecycle-il-operator/blob/fd4715c0e2978be8b3f3d0b7dbc33584ad390ff5/controllers/argocd/generate.go#L101) by the zLifecycle Operator.

### Local Development

Create a branch in the helm-charts repo and make changes to `terraform-config` chart.
Update your local zLifecycle Operator to point ArgoCD Applications to your branch of the helm chart, then redeploy your operator.

```
    Source: appv1.ApplicationSource{
        RepoURL:        helmChartsRepo,
        Path:           "charts/terraform-config",
        TargetRevision: "my-cool-branch",
        Helm: &appv1.ApplicationSourceHelm{
            Values: helmValues,
        },
    },
```