# Variables

## How they work in web

Since zlifecycle-web is a React app, environment variables are embedded during build time. In our infrastructure, build time corresponds to whenever the Dockerfile is built and tagged.

Whenever npm run build is used, according to the [CRA docs](https://create-react-app.dev/docs/adding-custom-environment-variables/), NODE_ENV will always be production and have the following priority for .env files:

```
npm run build: .env.production.local, .env.local, .env.production, .env
```

This ordering is important because in .env.production, variables that need to be altered at runtime in our containers are typically assigned a value with a prefix of \_\_DOCKER\_, e.g. REACT\_APP\_CUSTOMER\_NAME=\_\_DOCKER\_REACT\_APP\_CUSTOMER\_NAME\_\_. The reason for this is that there's a sed script being run in the Dockerfile that statically replaces the \_\_DOCKER\_REACT prefixed vars with their equivalent values from the environment.

Since we aren't currently running Docker locally, this doesn't happen locally, so you will need to provide the actual value in the corresponding development .env file.

## What does this mean for a developer?

### If you need a variable that has the same value in all environments

This is the simplest case. If you need this, just add it into the base .env file and reference it accordingly.  To use a value within the React app, make sure that it's prefixed REACT_APP.

### If you need a variable that has different values in an environment

If you need this, you need to do the following:

#### For production and sandbox

- Add in a REACT_APP_YOUR_VALUE=\_\_DOCKER\_REACT\_APP\_YOUR\_VALUE\_\_ to the .env.production configuration
- Make sure \_\_DOCKER\_REACT\_APP\_YOUR\_VALUE\_\_ has a value in the various helm environment deploy scripts

#### To see the value locally

- Add in a REACT_APP_YOUR_VALUE=VALUE to your various local .env files on your machine. Until we move to docker-compose, you'll need to manually do this.
