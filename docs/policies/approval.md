# Approval

If you wish for the Terraform plan to apply without requiring a manual approval step, you can do so by setting a property called `autoApprove`.

When you are [defining an environment](/define/define_environment), you can set the `autoApprove` flag to `true`, at the **environment** level or **component** level.

When the plan is run, the reconciler looks for this property which tells it whether to ask for explicit approval at the time of provisioning, teardown, or reconciling the environment - or - to automatically approve the plan.

This is an optional field, with the default value as `false`.

In the case that `autoApprove` flag was not specified or set to `false`, a prompt for approval will appear for each component being provisioned or destroyed. Details on how to approve are [here](/policies/manual_approval).

**NOTE:**
The `autoApprove` flag can be configured for all components in an environment at the spec level, and can also be specified at the component level, which will override the spec level configuration for that environment.

---
**Sample YAML**

```yaml
apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
  name: dev-checkout-sandbox
  namespace: zlifecycle
spec:
  teamName: checkout
  envName: sandbox  
  autoApprove: false # spec level
  components:
    - name: networking
      autoApprove: true # this will override the one at spec level
      type: terraform
      module:
        source: aws
        name: vpc
      variablesFile:
        source: "git@github.com:githubRepo.git" # Add your repo here
        path: "sandbox/tfvars/networking.tfvars" # Add your tfvars here
      outputs:
        - name: vpc_id
        - name: public_subnets
        - name: private_subnets
        - name: vpc_cidr_block
```
---
