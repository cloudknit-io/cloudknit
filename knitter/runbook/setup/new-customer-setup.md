# Setup a new Customer

## Overview

Steps to setup a new Customer

## When to use this runbook
This is to be used when you are setting up brand new customer

## Initial Steps Overview

1. [Setup by CompuZest](#setup-by-compuzest)
1. [Setup by Customer](#setup-by-customer)

## Detailed Steps

#### Setup by CompuZest
This account will be used by zLifecycle to Read and Write to various repos:

- Create a mailing group for `<client>-zlifecycle@compuzest.com` on G Suite
- Create new Github Service Account and register it under `<client>-zlifecycle@compuzest.com`, username should follow the format `<client>-zlifecycle`
- Generate Personal Access Token for the Zlifecycle Service Account to be used by the Operator (Check LastPass secret note: "zLifecycle - k8s secrets")
  - In the scope select all options for `repo` and `workflow`
- Generate SSH key for the Github Service Account to be used by the Operator (Check LastPass secret note: "zLifecycle - k8s secrets")
    ```shell script
    ssh-keygen -b 2048 -t rsa -f ~/.ssh/<client> -q -N "" -C "<client>-zlifecycle@compuzest.com"
    ```
After generating the SSH key make sure you add the public key to the Github `<client>-zlifecycle` service account by going to `Settings -> SSH and GPG keys`
- Save all the above secrets in LastPass under the svc account password entry in notes. See `compuzest@compuzest.com` entry as example 

#### Setup by Customer
- Add zLifecycle as an OAuth application by going `Repository Settings -> Developer Settings -> OAuth Apps -> New OAuth App`
    * Application name: `<client>-zlifecycle`
    * Homepage URL: `https://<client>.zlifecycle.com`
    * Application description (OPTIONAL): `zLifecycle instance for <client>`
    * Authorization callback URL: `https://<client>.zlifecycle.com/api/dex/callback`
- Generate a new client secret from the Application OAuth page
- Add Environment Config in `compuzest` zl env. See `https://github.com/compuzest-tech/zl-environment-config`
- Create a repo `zlifecycle-config` in customer github org and give `<client>-zlifecycle` svc account admin access to the repo
- Add Customer CR
- Add Team CRD in `Teams` folder
- Create team repo and give `<client>-zlifecycle` svc account admin access to the repo
- Add Team CRD and push


#### Setup AWS Creds
- Go to <client>.zlifecycle.com and login with creds and go to settings page and setup aws creds
