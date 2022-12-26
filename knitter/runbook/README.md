# Runbooks
Runbooks are a compilation of routine procedures with clear and complete steps for investigating, diagnosing, and fixing a defined problem.

* [Deploying Knitter with `kubectl`](#deploying-knitter-with-kubectl)

## Deploying Knitter with `kubectl`

Replace `[knitter-image-tag]` with the tag you're trying to deploy:

```bash
kubectl patch workflowtemplate audit-run-template --type='json' -p='[{"op": "replace", "path": "/spec/templates/0/script/image", "value":"413422438110.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-terraform:[knitter-image-tag]"}]' -n zlab-executor && \
kubectl patch workflowtemplate workflow-trigger-template --type='json' -p='[{"op": "replace", "path": "/spec/templates/2/script/image", "value":"413422438110.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-terraform:[knitter-image-tag]"}]' -n zlab-executor && \
kubectl patch workflowtemplate terraform-run-template --type='json' -p='[{"op": "replace", "path": "/spec/templates/0/script/image", "value":"413422438110.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-terraform:[knitter-image-tag]"}]' -n zlab-executor
```
