## Docker Compose setup

# Overview

Setup project locally for development using docker-compose

# When to use this runbook

If dev environment is not setup already with required environment variables.

# Steps

- clone [zlifecycle-web](https://github.com/CompuZest/zlifecycle-web)
- clone [zlifecycle-web-bff](https://github.com/CompuZest/zlifecycle-web-bff)
- create .env.development.docker file in ${ZLIFECYCLE_WEB_PROJECT_DIR}
```
HOST=0.0.0.0
PORT=3001

NODE_ENV=development
CHOKIDAR_USEPOLLING=true

REACT_APP_BASE_URL=http://localhost:8088
REACT_APP_STREAM_URL=http://localhost:8088
REACT_APP_AUTHORIZE_URL=http://localhost:8088/authorize

REACT_APP_SENTRY_ENVIRONMENT=development
```
- create .env.development.docker file in ${ZLIFECYCLE_BFF_PROJECT_DIR}
```
PORT=8088

COOKIE_SAME_SITE=lax
COOKIE_DOMAIN=localhost
COOKIE_SECRET=test

OPENID_ISSUER=http://host.docker.internal:8089/api/dex
OPENID_CLIENT_ID=<OPENID_CLIENT_ID>
OPENID_CLIENT_SECRET=<OPENID_CLIENT_SECRET>
OPENID_CALLBACK=http://localhost:8088/auth/callback

ARGO_WORKFLOW_API_URL=http://host.docker.internal:2746
ARGO_CD_API_URL=http://host.docker.internal:8089
ZLIFECYCLE_API_URL=http://host.docker.internal:4000
```

- run  ```ZLIFECYCLE_WEB_PROJECT_DIR=<path to web> ZLIFECYCLE_BFF_PROJECT_DIR=<path to bff> docker-compose up``` or ```docker-compose up -d``` for detached mode
    - for example, Ryan runs: ```ZLIFECYCLE_WEB_PROJECT_DIR=/Projects/ZLifeCycle/zlifecycle-web ZLIFECYCLE_BFF_PROJECT_DIR=/Projects/ZLifeCycle/zlifecycle-web-bff docker-compose up```
- visit [localhost:3001](http://localhost:3001)
