import helper from "../utils/helper";
import * as express from "express";
import { createClient } from "redis";
import EventEmitter = require("events");

const event = new EventEmitter();
let redisClient = null;

async function startRedis() {
  console.log("Starting Redis");
  redisClient = createClient({
    url: "redis://172.16.244.62:6379"
  });

  redisClient.on("error", (err) => console.log("Redis Client Error", err));

  await redisClient.connect();
  redisClient.SUBSCRIBE("test-channel", (message: any, channel: any) => {
    event.emit("stream", JSON.parse(message));
  });
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
      response.write(`event: ${stream.type}\ndata: ${JSON.stringify(stream.data)}\n\n`);
    }
  });

  request.on("close", () => {
    console.log(`Connection closed`);
  });
}

export function setUpSSE(router: express.Router) {
  router.get("/stream", eventsHandler);
}
