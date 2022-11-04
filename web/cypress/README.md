# Cypress

## Instructions

If you have the zlifecycle servers running via docker-compose, you can run
cypress locally to write e2e/integration tests by running the following command in the root folder:

```npm run cypress:open```

To run the tests headlessly against the docker-compose instances, you can run
the cypress tests headlessly by running the following command in the root
folder:

```ZLIFECYCLE_WEB_PROJECT_DIR=<YOUR PATH HERE> ZLIFECYCLE_BFF_PROJECT_DIR=<YOUR PATH HERE> docker-compose -f docker-compose.yml -f cypress/docker-compose.cypress.yml build```
```ZLIFECYCLE_WEB_PROJECT_DIR=<YOUR PATH HERE> ZLIFECYCLE_BFF_PROJECT_DIR=<YOUR PATH HERE> docker-compose -f docker-compose.yml -f cypress/docker-compose.cypress.yml up```

The tests will fail since it'll take a few minutes for the infrastructure to come up. Once the infrastructure is up, you can run the following:

```ZLIFECYCLE_WEB_PROJECT_DIR=<YOUR PATH HERE> ZLIFECYCLE_BFF_PROJECT_DIR=<YOUR PATH HERE> docker-compose -f docker-compose.yml -f cypress/docker-compose.cypress.yml run --rm tests```
