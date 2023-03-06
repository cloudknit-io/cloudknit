# Target base: installs dependencies
FROM node:18.14.2-alpine AS base

WORKDIR /var/www

COPY package.json package.json
COPY package-lock.json package-lock.json

ENV NODE_ENV production
# RUN npm ci --only=production // want to use this, but some types are in devDeps
RUN npm ci --production=false -f

FROM base as copy-app-stage

COPY . /var/www

# Target lint-stage: lints app
FROM copy-app-stage as lint-stage

# RUN npm run lint

# Target test-stage: runs tests
FROM lint-stage AS test-stage
RUN CI=true npm run test

# TODO: Clean up this stage with .dockerignore
# Target build-stage: builds app
FROM test-stage AS build-stage

RUN npm run build:production

# TODO: Look at nginx.conf more closely
# TODO: Get rid of bash dependency; convert to shell
# TODO: Look at env.sh perms (ie is 755 correct?)
# Target production: serves app
FROM nginx:stable-alpine as production

RUN apk add bash

WORKDIR /usr/share/nginx/html

RUN rm -rf ./*

COPY --from=build-stage /var/www/conf/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage --chmod=0755 /var/www/env.sh env.sh
COPY --from=build-stage /var/www/.env.production .env.production
COPY --from=build-stage /var/www/build .

EXPOSE 80

# Containers run nginx with global directives and daemon off
ENTRYPOINT ["/bin/bash", "-c", ". /usr/share/nginx/html/env.sh && exec nginx -g 'daemon off;'"]
