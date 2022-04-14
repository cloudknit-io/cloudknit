# Setup a new Company

Steps to setup a new Company

### Setup by zLifecycle

#### Provision new company environment

- Add Environment Config in `compuzest` zl env. See `https://github.com/compuzest-tech/zl-environment-config`
    * Use dummy secrets that you don't have yet from the company (like oAuth secrets)

### Setup by Company

#### Create zLifecycle config repos

- Create a repo `zlifecycle-config` in company github org and give `<client>-zlifecycle` svc account admin access to the repo
- Add Team YAML in `Teams` folder
- Create team repo 
- Push code to Github

#### Create OAuth application for login to zLifecycle

- Add zLifecycle as an OAuth application by going `Repository Settings -> Developer Settings -> OAuth Apps -> New OAuth App`
    * Application name: `<client>-zlifecycle`
    * Homepage URL: `https://<client>.app.zlifecycle.com`
    * Application description (OPTIONAL): `zLifecycle instance for <client>`
    * Authorization callback URL: `https://<client>.app.zlifecycle.com/api/dex/callback`
- Generate a new client secret from the Application OAuth page
- Enter secret on zLifecycle UI

#### Give access to zLifecycle config & terraform modules repo

Use one of the options below:

##### Using Github Service Account

- Give `<client>-zlifecycle` Github svc account admin access to the repo
- Or give readonly access and create the webhook manually 

##### Using zLifecycle Github App

Install zLifecycle Github App to repositories so that zLifecycle can read config and create webhook

- Go to [zLifecycle Github App](https://github.com/apps/zlifecycle)
- Install App to your Github Org & give access to zlifecycle config & terraform module repos. 
    * This will give read only access & access to create Webhook

### Setup by zLifecycle

#### Update secrets from Company

- Update environment with Company secrets
