## Local BFF Setup

# Overview
Setup BFF locally for development.

# When to use this runbook
If dev environment is not setup already with required environment variables.

# Steps
-   clone repository from git
-   install all required libraries: `npm install`
-   Port forward following services using command `kubectl port-forward {service-name} {localPort:envPort} -n {namespace}`
    -   service-name: `service/argocd-{organization}-server`, localPort:envPort: `8081:80` namespace: `{organisation}-system`
    -   service-name: `service/argo-workflow-server`, localPort:envPort: `2746:2746`,  namespace: `{organisation}-executor`
    -   service-name: `svc/zlifecycle-api`, localPort:envPort: `4000:80`, namespace: `{organisation}-system`
    -   service-name: `svc/zlifecycle-state-manager`, localPort:envPort: `5000:80`, namespace: `{organisation}-system`
    -   service-name: `svc/event-service`, localPort:envPort: `8082:8081`, namespace: `{organisation}-system`
    -   service-name: `svc/event-service`, localPort:envPort: `8083:8082`, namespace: `{organisation}-system`
-   Update/Add env file with following inputs
    - PORT=8080
    - SITE_URL=http://localhost:3000
    - ARGO_WORKFLOW_API_URL=http://localhost:2746
    - ARGO_CD_API_URL=http://localhost:8081
    - ZLIFECYCLE_API_URL=http://localhost:4000
    - ZLIFECYCLE_STATE_MANAGER_URL=http://localhost:5000
    - ZLIFECYCLE_EVENT_API_URL=http://localhost
-   start project: `npm run start`
-   Or add the above given variables prepended with export command and add them to a shell file and execute that file to start the project.

    **Make sure you are running node version 14+**

    **Common Errors:**
    - the error `Error: Cannot find module 'jose/jwe/compact/encrypt'` might mean you are not running Node version 14.5
    - the error `.keys required.` might mean env variables are not set
