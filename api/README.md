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
1. Run `docker-compose down -v --remove-orphans` to kill the MySQL container

### RDS

1. Connect to [VPN](https://github.com/CompuZest/engineering/blob/main/docs/onboarding.md).
1. Set the proper credentials in `run.sh`
1. You can get RDS credentials by viewing the helm chart of `zlifecycle-api` pod

### Generating Migrations

1. Ensure you're on an up to date schema version by running `./run.sh`
1. Make changes to the required entity
1. Generate Migrations: `npm run typeorm migration:generate ./src/typeorm/migrations/[Name of migration]`
1. Migrations will be generated as TypeScript. Unfortunately, I've not been able to them to run successfully as TypeScript so you'll have to convert them to JavaScript. This is easy to do. Look [here](src/typeorm/migrations/1673901210875-TeamEstimatedCost.js) for an example.

### Applying Migrations

You'll need a working `.env`:

1. `cp .env.example .env` - local
1. `cp .env.example .env.dev` - dev
1. `cp .env.example .env.prod` - prod

> Note: Look in AWS, 1Password, k9s for values

* Locally
   * migrations are applied automatically when running on your machine
   * You can find that configuration [here](src/typeorm/index.ts)
* Dev / Prod
   1. connect to VPN
   1. `npm run typeorm:[dev|prod] migrate:run`
