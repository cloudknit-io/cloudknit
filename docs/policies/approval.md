# Approval

When you are creating an environment YAML file, you have an optional property called **autoApprove** which can be added at **environment** level or **component** level.

Our reconciler, looks for this property which tells it whether to ask for explicit approval at the time of provisioning, teardown, or reconciling the environment.

**NOTE:**

* Approval prompt will only appear if **autoApprove** is `false`
* The default value of **autoApprove** is `false`
* Component level **autoApprove** flag overrides spec level **autoApprove**

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
