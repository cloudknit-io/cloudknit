# Setup steps 

## Create OAuth application for login to zLifecycle

- Add zLifecycle as an OAuth application by going `Repository Settings -> Developer Settings -> OAuth Apps -> New OAuth App`
    * Application name: `<client>-zlifecycle`
    * Homepage URL: `https://<client>.app.zlifecycle.com`
    * Application description (OPTIONAL): `zLifecycle instance for <client>`
    * Authorization callback URL: `https://<client>.app.zlifecycle.com/api/dex/callback`
- Generate a new client secret from the Application OAuth page
- Enter secret on zLifecycle UI Settings page

## Create zLifecycle config repos

- Create a repo `zlifecycle-config` in your github org
- Onboard a Team using [this](https://docs.zlifecycle.com/getting_started/onboard_team/)

## Give access to zLifecycle config & terraform modules repo

Install zLifecycle Github App using steps [here](https://docs.zlifecycle.com/settings/zlifecycle_app_installation/)
