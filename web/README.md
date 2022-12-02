# zLifecycle

# Prerequisites
- [Boostrap AWS & K8s](https://github.com/cloudknit-io/cloudknit/tree/main/web/runbook/bootstrap/set-up-aws-k8s.md)

## Starting project
- [Setting up Web](https://github.com/cloudknit-io/cloudknit/tree/main/web/runbook/bootstrap/web-setup.md)

## ZLifecycle Local environment
- Port forward ArgoCD at port `8081`, start the `zlifecycle-web-bff` application at port `8080` which is its default and run `cp .env.development .env.development.local`
```bash
kubectl port-forward service/argocd-server 8081:80 -n argocd
```
- Port forward Argo Workflow using following command
```bash
kubectl port-forward service/argo-workflow-server 2746:2746 -n argocd
```

## Node and npm
-   [NVM](https://github.com/creationix/nvm) is used to specify Node version
-   On some hardware, node version 12.3.0 is required, on others 14+ is possible (node-sass has install requirements)
-   when adding a new package exact version should be used and dependencies should be fixed with `npm audit fix`
-   to update all package versions: `npx npm-check-updates -u` or `ncu`
-   if there is change `http/https` in `package-lock.json` when you run `npm install`, then you should do following: `npm cache clean --force`, `rm -rf node_modules` and then `npm install`. This is a known issue with npm
-   update your npm version to latest version 

## Useful scripts
-   run tests: `npm test` (update snapshots: `npm test -- -u`)
-   check for lint errors: `npm run lint`
-   create build version: `npm build`

## Start project
-   run project locally: `npm run start:local`
-   run project on specific environment: `npm run start:[env]` (where `env` can be `local`, `dev`, `production` etc.)

## Build versions
-   create a local build version: `npm run build`
-   create a build version: `npm run build:[env]` (where `env` can be `dev`, `production` etc.)

## Libraries used
-   [Formik](https://jaredpalmer.com/formik)
-   [React router](https://reacttraining.com/react-router/)
-   [JWT Decode](https://github.com/auth0/jwt-decode)
-   [Axios](https://github.com/axios/axios)
-   [Monaco React](https://github.com/suren-atoyan/monaco-react)
-   [Argo CD - Embedded](https://github.com/argoproj/argo-cd)

## Prettier recommended setup

### Webstorm
-   Install Prettier formatter extension, check [this](https://plugins.jetbrains.com/plugin/10456-prettier) for more info

### Visual Studio Code
-   Install Prettier-Code formatter extension, check [this](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) for more info
-   Install ESLint extension, check [this](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) for more info
-   In root of project create .vscode folder
-   Inside .vscode folder create settings.json
-   Open settings.json and paste following code from [vscode-settings.json](./assets/vscode-settings.json)

## Useful links
-   [Typescript](https://www.typescriptlang.org/)
-   [About Prettier](https://www.jetbrains.com/help/idea/prettier.html)

## Testing
-   Run only one test file -  `npm run test -- --watch -f './testname.test.js' --coverage=false --verbose=false`

### Automating tests
-   Add `data-testid` on user interactive components so that finding feature during automated testing is easier.

## Troubleshooting
If in terraform a helm resource gets orphaned, trying to re-apply that helm release will create the error `Error: cannot re-use a name that is still in use`. Check that every resource including the helm generated secret of `sh.helm.release.v1.zlifecycle-web.v1` has been destroyed in the cluster and try again
