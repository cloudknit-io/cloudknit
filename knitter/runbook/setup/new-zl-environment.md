# Steps for new zLifecycle Environment

1. [Setup by CompuZest](#setup-by-compuzest)

#### Create Github Service Account (To be done once per zlEnvironment like dev,app)

Steps to follow:

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
