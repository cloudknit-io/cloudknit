# Setup steps 

Setup steps for zLifecycle

## Create OAuth application for login to zLifecycle

- Add zLifecycle as an OAuth application by going `Repository Settings -> Developer Settings -> OAuth Apps -> New OAuth App`
    * Application name: `<client>-zlifecycle`
    * Homepage URL: `https://<client>.app.zlifecycle.com`
    * Application description (OPTIONAL): `zLifecycle instance for <client>`
    * Authorization callback URL: `https://<client>.app.zlifecycle.com/api/dex/callback`
- Generate a new client secret from the Application OAuth page
- Enter secret on zLifecycle UI

## Create zLifecycle config repos

- Create a repo `zlifecycle-config` in your github org
- Add Team YAML in `Teams` folder
- Create team repo 
- Push code to Github

## Give access to zLifecycle config & terraform modules repo

Install zLifecycle Github App iusing steps [here](https://docs.zlifecycle.com/settings/zlifecycle_app_installation/)
