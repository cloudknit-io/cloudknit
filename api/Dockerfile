FROM node:16.18-alpine3.16 AS build

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

FROM node:16.18-alpine3.16

WORKDIR /app

COPY package.json .
RUN npm install --production

COPY --from=build ./dist/ ./dist/

EXPOSE 3000

CMD [ "npm", "run", "start:prod" ]
