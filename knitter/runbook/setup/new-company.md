# Setup a new Company

Steps to setup a new Company

### Setup by Company

#### Create OAuth application for login to zLifecycle (To be done per Company)

- Add zLifecycle as an OAuth application by going `Repository Settings -> Developer Settings -> OAuth Apps -> New OAuth App`
    * Application name: `<client>-zlifecycle`
    * Homepage URL: `https://<client>.zlifecycle.com`
    * Application description (OPTIONAL): `zLifecycle instance for <client>`
    * Authorization callback URL: `https://<client>.zlifecycle.com/api/dex/callback`
- Generate a new client secret from the Application OAuth page
- Add Environment Config in `compuzest` zl env. See `https://github.com/compuzest-tech/zl-environment-config`
- Create a repo `zlifecycle-config` in company github org and give `<client>-zlifecycle` svc account admin access to the repo
- Add Team YAML in `Teams` folder
- Create team repo and give `<client>-zlifecycle` svc account admin access to the repo
- Push code to Github

#### Install zLifecycle Github App

We need to install zLifecycle Github App to repositories so that zLifecycle can read from them and create webhooks

- Go to [zLifecycle Github App](https://github.com/organizations/zLifecycle/settings/apps/zlifecycle/installations)
- Install App to Companies Org & give access to config repos (like `zlifecycle-config` & `checkout-team-config`)
