# Target base: initializes
FROM node:lts-alpine3.16 AS base

WORKDIR /var/www

# Target install-deps-stage: installs dependencies
FROM base AS install-deps-stage

COPY ["package.json", "package.json"]
COPY ["package-lock.json", "package-lock.json"]

ENV NODE_ENV production
RUN npm ci --only=production

# TODO: Clean up this stage with .dockerignore
# Target build-stage: builds app
FROM base AS build-stage

WORKDIR /var/www

COPY --from=install-deps-stage /var/www/node_modules /var/www/node_modules
COPY . /var/www

RUN npm run build

FROM base AS production

COPY --from=install-deps-stage /var/www/node_modules /var/www/node_modules
COPY --from=install-deps-stage /var/www/package.json /var/www/package.json
COPY --from=build-stage /var/www/dist /var/www/dist

EXPOSE 8080

ENTRYPOINT ["npm", "run", "server"]
