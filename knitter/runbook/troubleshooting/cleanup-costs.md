# Cleanup Dangling Costs for Teams/Envs

## Overview

Steps to fix an issue that happens periodically. The costs at Team/Env level are incorrect.

## Initial Steps Overview

* Port forward API service by using below command.

```bash
kubectl port-forward svc/zlifecycle-api 4000:80 -n zmart-system
```

* Cleanup costs

```bash
curl -X 'GET' \
  'http://localhost:4000/costing/api/v1/execute/Update%20components%20set%20isDeleted%20%3D%201%20where%20team_name%20%3D%20%%27%20and%20environment_name%20%3D%20%27---envName---%27' \
  -H 'accept: */*'*

  ```
