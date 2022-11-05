# zLifecycle Internal API

## Description

[Nest](https://github.com/nestjs/nest) framework TypeScript starter repository.

## Installation

```bash
$ npm install
```

## Running Locally

1. `cp run.sh.example run.sh`
1. Configure your `run.sh` with the appropriate organizations

### Local MySQL Database

1. Requires [Docker](https://docs.docker.com/desktop/#download-and-install)
1. `./run.sh` will automatically start a local MySQL instance
1. Run `docker-compose down -v` to kill the MySQL container

### RDS

1. Connect to [VPN](https://github.com/CompuZest/engineering/blob/main/docs/onboarding.md).
1. Set the proper credentials in `run.sh`
1. You can get RDS credentials by viewing the helm chart of `zlifecycle-api` pod
