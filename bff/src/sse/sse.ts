import helper from "../utils/helper";
import * as express from "express";
import { createClient } from "redis";
import EventEmitter = require("events");
import logger from "../utils/logger";
import config from "../config";

const event = new EventEmitter();
let redisClient = null;
const apiStreamChannel = 'api-stream-channel';

async function startRedis() {
  try {
    redisClient = createClient({
      url: config.redis.url,
      password: config.redis.password,
    });

    redisClient.on("error", (err) => console.log("Redis Client Error", err));

    await redisClient.connect();

    if (redisClient.isReady) {
      redisClient.SUBSCRIBE(apiStreamChannel, (message: any, channel: any) => {
        event.emit("stream", JSON.parse(message));
      });
    }
  } catch (err) {
    redisClient = null;
    logger.error('Redis failed to connect', err);
  }
}

async function eventsHandler(request: any, response: any, next) {
  const reqUser = helper.userFromReq(request);
  if (!reqUser) {
    response.status(400).send({
      unauthorized: "User is not authorized",
    });
    return;
  }

  if (!redisClient) {
    await startRedis();
  }

  const headers = {
    "Content-Type": "text/event-stream",
    Connection: "keep-alive",
    "Cache-Control": "no-cache",
  };

  response.writeHead(200, headers);

  event.on("stream", (stream: { data: any; type: string }) => {
    if (reqUser.organizations.some((org) => org.id === stream.data.orgId)) {
      response.write(
        `event: ${stream.type}\ndata: ${JSON.stringify(stream.data)}\n\n`
      );
    }
  });

  request.on("close", () => {
    console.log(`Connection closed`);
  });
}

export function setUpSSE(router: express.Router) {
  router.get("/session/stream", eventsHandler);
}
