FROM node:14 AS build

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

FROM node:14

WORKDIR /app

COPY package.json .
RUN npm install --production

COPY --from=build ./dist/ ./dist/
COPY --from=build ./ormconfig.js .

EXPOSE 3000

CMD [ "npm", "run", "start:prod" ]
