## zLifecycle Auth and Proxy server

### How to start

-   Follow the prerequisites given in [zlifecycle-web](https://github.com/CompuZest/zlifecycle-web/blob/main/README.md) repo if not already done.

#### Set up
[Setting up locally](https://github.com/CompuZest/zlifecycle-web-bff/blob/main/runbook/setting-up-bff.md)
#### Install

    npm install
#### Build

    npm run build

#### Configuration put in .env.local

check secret for web bff demo or sandbox, this has the env vars you need to export with the correct auth keys


curl https://dev-04d2288z.us.auth0.com/.well-known/jwks.json -i -H "Accept: /" -H "Origin: http://localhost:8080"
