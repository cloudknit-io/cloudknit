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
    > Note: zbank and local already exist. Copy and paste as needed.

### Local MySQL Database

1. Requires [Docker](https://docs.docker.com/desktop/#download-and-install)
1. `sh run.sh local`
    * This starts the API via `npm run start:debug` and a MySQL instance via `docker-compose`
1. Run `docker-compose down -v` to kill the MySQL container

### RDS

1. Connect to [VPN](https://github.com/CompuZest/engineering/blob/main/docs/onboarding.md).
1. Set the proper credentials in `run.sh` for the org you want to use
1. You can get RDS credentials by viewing the helm chart of `zlifecycle-api` pod
1. `sh run.sh zbank`
    > Note: `run.sh` uses named AWS profiles such as `zl-dev`. It's advisable to follow this convention but not required. Modify `run.sh` to suit your needs.

## Test

```bash
# unit tests
$ npm run test

# e2e tests
$ npm run test:e2e

# test coverage
$ npm run test:cov
```

## Support

Nest is an MIT-licensed open source project. It can grow thanks to the sponsors and support by the amazing backers. If you'd like to join them, please [read more here](https://docs.nestjs.com/support).

## Stay in touch

- Author - [Kamil Myśliwiec](https://kamilmysliwiec.com)
- Website - [https://nestjs.com](https://nestjs.com/)
- Twitter - [@nestframework](https://twitter.com/nestframework)

## License

Nest is [MIT licensed](LICENSE).
