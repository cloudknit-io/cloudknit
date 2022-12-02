## Local Web Setup 

# Overview
Setup project locally for development.

# When to use this runbook
If dev environment is not setup already with required environment variables.

# Steps
-   clone repository from git
-   install all required libraries: `npm install`
-   Update local env file (.env.development.local) with following inputs:
    - HOST=localhost
    - PORT=3000
    - REACT_APP_BASE_URL=/
    - REACT_APP_BASE_URL=http://localhost:8080
    - REACT_APP_STREAM_URL=http://localhost:8080
    - REACT_APP_AUTHORIZE_URL=http://localhost:8080/authorize
    - REACT_APP_SENTRY_ENVIRONMENT=local
    - REACT_APP_ENABLED_FEATURE_FLAGS=```PAGE_DASHBOARD|PAGE_BUILDER|PAGE_APPLICATIONS|TAB_DETAILED_LOGS|TAB_STATE_FILE|HARD_SYNC|DIFF_CHECKER|VISUALIZATION|BLUE_GREEN_DEPLOYMENT|TERM_AGREEMENT|QUICK_START```
-   start project: `npm run start:local`
-   Or add the above given variables prepended with export command and add them to a shell file and execute that file to start the project.
-   if not automatically, visit [localhost:3000](http://localhost:3000)
-   if you want to try localhost with https: `HTTPS=true npm run start:local`
-   To Login you first have to clone, setup and run [zlifecycle-web-bff](https://github.com/cloudknit-io/cloudknit/tree/main/bff/README.md) first.
