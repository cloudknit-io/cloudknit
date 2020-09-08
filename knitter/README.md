# infra-deploy-bootstrap
Bootstrap for infra deployments

## Bootstrap Terraform Shared State

Run `tfstate` terraform to provision S3 bucket and Dynamo DB table that will be used for Terraform Shared State.

```bash
cd tfstate
terraform init
terraform plan
terraform apply
```
